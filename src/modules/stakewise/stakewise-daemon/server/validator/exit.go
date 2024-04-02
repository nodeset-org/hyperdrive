package swvalidator

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	duserver "github.com/nodeset-org/hyperdrive/daemon-utils/server"
	api "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/node/validator"
	"github.com/rocket-pool/node-manager-core/utils/input"
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
	duserver.RegisterQuerylessGet[*validatorExitContext, api.ValidatorExitData](
		router, "exit", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
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

func (c *validatorExitContext) PrepareData(data *api.ValidatorExitData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	bc := sp.GetBeaconClient()
	w := sp.GetWallet()
	ctx := c.handler.ctx

	if len(c.pubkeys) == 0 {
		return types.ResponseStatus_Success, nil
	}

	// Requirements
	err := sp.RequireBeaconClientSynced(ctx)
	if err != nil {
		return types.ResponseStatus_ClientsNotSynced, err
	}

	// Load the keys
	keys := make([]*eth2types.BLSPrivateKey, len(c.pubkeys))
	for i, pubkey := range c.pubkeys {
		key, err := w.GetPrivateKeyForPubkey(pubkey)
		if err != nil {
			return types.ResponseStatus_Error, err
		}
		keys[i] = key
	}

	// Get the epoch of the chain head if needed
	if !c.isEpochSet {
		head, err := bc.GetBeaconHead(ctx)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting beacon head: %w", err)
		}
		c.epoch = head.Epoch
	}

	// Get the BlsToExecutionChange signature domain
	signatureDomain, err := bc.GetDomainData(ctx, eth2types.DomainVoluntaryExit[:], c.epoch, false)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting Beacon domain data: %w", err)
	}

	// Get the statuses (indices) of each validator
	statuses, err := bc.GetValidatorStatuses(ctx, c.pubkeys, nil)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting validator indices: %w", err)
	}

	// Get the signatures
	data.ExitInfos = make([]api.ValidatorExitInfo, len(keys))
	for i, key := range keys {
		// Get signed voluntary exit message
		pubkey := c.pubkeys[i]
		index := statuses[pubkey].Index

		signature, err := validator.GetSignedExitMessage(key, index, c.epoch, signatureDomain)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting exit message signature for validator %s: %w", pubkey.Hex(), err)
		}
		indexUint, _ := strconv.ParseUint(index, 10, 64)

		data.ExitInfos[i] = api.ValidatorExitInfo{
			Pubkey:    pubkey,
			Index:     indexUint,
			Signature: signature,
		}
		if !c.noBroadcast {
			err = bc.ExitValidator(ctx, index, c.epoch, signature)
			if err != nil {
				return types.ResponseStatus_Error, fmt.Errorf("error exiting validator %s: %w", pubkey.Hex(), err)
			}
		}
	}

	return types.ResponseStatus_Success, nil
}
