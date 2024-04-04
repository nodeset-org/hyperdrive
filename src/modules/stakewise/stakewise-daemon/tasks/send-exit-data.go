package swtasks

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/log"
	"github.com/rocket-pool/node-manager-core/node/validator"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	PubkeyKey string = "pubkey"
)

// Send exit data task
type SendExitData struct {
	logger *log.Logger
	ctx    context.Context
	sp     *swcommon.StakewiseServiceProvider
	w      *swcommon.Wallet
	ns     *swcommon.NodesetClient
	bc     beacon.IBeaconClient
}

// Create Exit data task
func NewSendExitData(ctx context.Context, sp *swcommon.StakewiseServiceProvider, logger *log.Logger) *SendExitData {
	return &SendExitData{
		logger: logger,
		ctx:    ctx,
		sp:     sp,
		w:      sp.GetWallet(),
		ns:     sp.GetNodesetClient(),
		bc:     sp.GetBeaconClient(),
	}
}

// Update Exit data
func (t *SendExitData) Run() error {
	t.logger.Info("Checking for missing signed exit data...")

	// Get registered validators
	resp, err := t.ns.GetRegisteredValidators(t.ctx)
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}
	for _, status := range resp {
		t.logger.Debug(
			"Retrieved registered validator",
			slog.String(PubkeyKey, status.Pubkey.HexWithPrefix()),
			slog.Bool("uploaded", status.ExitMessageUploaded),
		)
	}

	// Check for any that are missing signed exits
	missingExitPubkeys := []beacon.ValidatorPubkey{}
	for _, v := range resp {
		if v.ExitMessageUploaded {
			continue
		}
		missingExitPubkeys = append(missingExitPubkeys, v.Pubkey)
		t.logger.Info("Validator is missing a signed exit message.", slog.String(PubkeyKey, v.Pubkey.HexWithPrefix()))
	}
	if len(missingExitPubkeys) == 0 {
		return nil
	}

	// Get statuses for validators with missing exits
	statuses, err := t.bc.GetValidatorStatuses(t.ctx, missingExitPubkeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}

	// Get beacon head and domain data
	head, err := t.bc.GetBeaconHead(t.ctx)
	if err != nil {
		return fmt.Errorf("error getting beacon head: %w", err)
	}
	epoch := head.Epoch
	signatureDomain, err := t.bc.GetDomainData(t.ctx, eth2types.DomainVoluntaryExit[:], epoch, false)
	if err != nil {
		return fmt.Errorf("error getting domain data: %w", err)
	}

	// Get signed exit messages
	exitData := []swcommon.ExitData{}
	for _, pubkey := range missingExitPubkeys {
		t.logger.Warn("!!! pubkey: %v", pubkey)

		key, err := t.w.GetPrivateKeyForPubkey(pubkey)
		if err != nil {
			// Print message and continue because we don't want to stop the loop
			t.logger.Warn("Error getting private key", slog.String(PubkeyKey, pubkey.HexWithPrefix()), log.Err(err))
			continue
		}
		if key == nil {
			t.logger.Warn("Private key is nil", slog.String(PubkeyKey, pubkey.HexWithPrefix()))
			continue
		}
		index := statuses[pubkey].Index
		if index == "" {
			t.logger.Warn("Validator index is empty", slog.String(PubkeyKey, pubkey.HexWithPrefix()))
			continue
		}
		t.logger.Warn("!!! key: %v", key)

		signature, err := validator.GetSignedExitMessage(key, index, epoch, signatureDomain)
		if err != nil {
			// Print message and continue because we don't want to stop the loop
			// Index might not be ready
			t.logger.Warn("Error getting signed exit message", slog.String(PubkeyKey, pubkey.HexWithPrefix()), log.Err(err))
			continue
		}
		exitData = append(exitData, swcommon.ExitData{
			Pubkey: pubkey.HexWithPrefix(),
			ExitMessage: swcommon.ExitMessage{
				Message: swcommon.ExitMessageDetails{
					Epoch:          strconv.FormatUint(epoch, 10),
					ValidatorIndex: index,
				},
				Signature: signature.HexWithPrefix(),
			},
		})
	}

	// Upload the messages to Nodeset
	if len(exitData) > 0 {
		t.logger.Warn("!!! SEND EXIT DATA")
		_, err := t.ns.UploadSignedExitData(t.ctx, exitData)
		if err != nil {
			return fmt.Errorf("error uploading signed exit messages to NodeSet: %w", err)
		}

		pubkeys := []any{}
		for _, d := range exitData {
			pubkeys = append(pubkeys, slog.String(PubkeyKey, d.Pubkey))
		}
		t.logger.Info("Uploaded exit messages", pubkeys...)
	}

	return nil
}
