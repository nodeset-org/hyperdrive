package swvalidator

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/daemon-utils/validator/utils"
	api "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	pubkeyLimit int = 100000 // Basically no limit
)

// ===============
// === Factory ===
// ===============

type validatorExitContextFactory struct {
	handler *ValidatorHandler
}

func (f *validatorExitContextFactory) Create(args url.Values) (*validatorExitContext, error) {
	c := &validatorExitContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateOptionalArg("epoch", args, input.ValidateUint, &c.epoch, &c.isEpochSet),
		server.ValidateArgBatch("pubkeys", args, pubkeyLimit, input.ValidatePubkey, &c.pubkeys),
		server.ValidateArg("no-broadcast", args, input.ValidateBool, &c.noBroadcast),
	}
	return c, errors.Join(inputErrs...)
}

func (f *validatorExitContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*validatorExitContext, api.ValidatorExitData](
		router, "exit", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type validatorExitContext struct {
	handler     *ValidatorHandler
	epoch       uint64
	isEpochSet  bool
	pubkeys     []beacon.ValidatorPubkey
	noBroadcast bool
}

func (c *validatorExitContext) PrepareData(data *api.ValidatorExitData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	bc := sp.GetBeaconClient()
	w := sp.GetWallet()

	if len(c.pubkeys) == 0 {
		return nil
	}

	// Requirements
	err := sp.RequireBeaconClientSynced(context.Background())
	if err != nil {
		return err
	}

	// Load the keys
	keys := make([]*eth2types.BLSPrivateKey, len(c.pubkeys))
	for i, pubkey := range c.pubkeys {
		key, err := w.GetPrivateKeyForPubkey(pubkey)
		if err != nil {
			return err
		}
		keys[i] = key
	}

	// Get the epoch of the chain head if needed
	if !c.isEpochSet {
		head, err := bc.GetBeaconHead(context.Background())
		if err != nil {
			return fmt.Errorf("error getting beacon head: %w", err)
		}
		c.epoch = head.Epoch
	}

	// Get the BlsToExecutionChange signature domain
	signatureDomain, err := bc.GetDomainData(context.Background(), eth2types.DomainVoluntaryExit[:], c.epoch, false)
	if err != nil {
		return fmt.Errorf("error getting Beacon domain data: %w", err)
	}
	// Get the statuses (indices) of each validator
	statuses, err := bc.GetValidatorStatuses(context.Background(), c.pubkeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator indices: %w", err)
	}
	// Get the signatures
	data.ExitInfos = make([]api.ValidatorExitInfo, len(keys))
	for i, key := range keys {
		// Get signed voluntary exit message
		pubkey := c.pubkeys[i]
		index := statuses[pubkey].Index

		signature, err := utils.GetSignedExitMessage(key, index, c.epoch, signatureDomain)
		if err != nil {
			return fmt.Errorf("error getting exit message signature for validator %s: %w", pubkey.Hex(), err)
		}
		indexUint, _ := strconv.ParseUint(index, 10, 64)

		data.ExitInfos[i] = api.ValidatorExitInfo{
			Pubkey:    pubkey,
			Index:     indexUint,
			Signature: signature,
		}
		if !c.noBroadcast {
			err = bc.ExitValidator(context.Background(), index, c.epoch, signature)
			if err != nil {
				return fmt.Errorf("error exiting validator %s: %w", pubkey.Hex(), err)
			}
		}
	}

	return nil
}
