package swtasks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/validator/utils"
	swcommon "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common"
	"github.com/nodeset-org/hyperdrive/shared/utils/log"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

// Send exit data task
type SendExitData struct {
	sp *swcommon.StakewiseServiceProvider

	log log.ColorLogger
}

// Create Exit data task
func NewSendExitData(sp *swcommon.StakewiseServiceProvider, logger log.ColorLogger) *SendExitData {
	return &SendExitData{
		sp:  sp,
		log: logger,
	}
}

// Update Exit data
func (t *SendExitData) Run() error {
	t.log.Println("Checking for missing signed exit data...")

	// Get services
	w := t.sp.GetWallet()
	ns := t.sp.GetNodesetClient()
	bc := t.sp.GetBeaconClient()

	// Get registered validators
	resp, err := ns.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}

	// Check for any that are missing signed exits
	missingExitPubkeys := []beacon.ValidatorPubkey{}
	for _, v := range resp {
		if v.Uploaded {
			continue
		}
		missingExitPubkeys = append(missingExitPubkeys, v.Pubkey)
		fmt.Printf("Validator %v is missing a signed exit message.\n", v.Pubkey)
	}
	if len(missingExitPubkeys) == 0 {
		return nil
	}

	// Get statuses for validators with missing exits
	statuses, err := bc.GetValidatorStatuses(context.Background(), missingExitPubkeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}

	// Get beacon head and domain data
	head, err := bc.GetBeaconHead(context.Background())
	if err != nil {
		return fmt.Errorf("error getting beacon head: %w", err)
	}
	epoch := head.Epoch
	signatureDomain, err := bc.GetDomainData(context.Background(), eth2types.DomainVoluntaryExit[:], epoch, false)
	if err != nil {
		return fmt.Errorf("error getting domain data: %w", err)
	}

	// Get signed exit messages
	exitData := []swcommon.ExitData{}
	for _, pubkey := range missingExitPubkeys {
		key, err := w.GetPrivateKeyForPubkey(pubkey)
		if err != nil {
			// Print message and continue because we don't want to stop the loop
			fmt.Printf("WARNING: error getting private key for pubkey %s: %s\n", pubkey.HexWithPrefix(), err.Error())
			continue
		}
		index := statuses[pubkey].Index
		signature, err := utils.GetSignedExitMessage(key, index, epoch, signatureDomain)
		if err != nil {
			// Print message and continue because we don't want to stop the loop
			// Index might not be ready
			fmt.Printf("WARNING: error getting signed exit message for pubkey %s: %s", pubkey.HexWithPrefix(), err.Error())
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
		_, err := ns.UploadSignedExitData(exitData)
		if err != nil {
			return fmt.Errorf("error uploading signed exit messages to NodeSet: %w", err)
		}

		fmt.Println("Registered validators:")
		for _, d := range exitData {
			fmt.Printf("\t%s\n", d.Pubkey)
		}
	}

	return nil
}
