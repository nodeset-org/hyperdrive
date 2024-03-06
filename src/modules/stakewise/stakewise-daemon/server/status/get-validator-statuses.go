package swstatus

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/nodeset-org/eth-utils/beacon"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
)

// ===============
// === Factory ===
// ===============

type statusGetValidatorsStatusesContextFactory struct {
	handler *StatusHandler
}

func (f *statusGetValidatorsStatusesContextFactory) Create(args url.Values) (*statusGetValidatorsStatusesContext, error) {
	c := &statusGetValidatorsStatusesContext{
		handler: f.handler,
	}
	inputErrs := []error{}
	return c, errors.Join(inputErrs...)
}

func (f *statusGetValidatorsStatusesContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*statusGetValidatorsStatusesContext, swapi.ValidatorStatusData](
		router, "status", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type statusGetValidatorsStatusesContext struct {
	handler *StatusHandler
}

func (c *statusGetValidatorsStatusesContext) PrepareData(data *swapi.ValidatorStatusData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	bc := sp.GetBeaconClient()
	w := sp.GetWallet()
	nc := sp.GetNodesetClient()
	registeredPubkeys, err := nc.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}
	privateKeys, err := w.GetAllPrivateKeys()
	if err != nil {
		return fmt.Errorf("error getting private keys: %w", err)
	}
	publicKeys, err := w.DerivePubKeys(privateKeys)
	if err != nil {
		return fmt.Errorf("error getting public keys: %w", err)
	}
	statuses, err := bc.GetValidatorStatuses(context.Background(), publicKeys, nil)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}

	validatorStatuses := make(map[beacon.ValidatorPubkey]swapi.ValidatorStatus)
	for _, pubKey := range publicKeys {
		switch {
		case IsWithdrawalDone(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.WithdrawalDone
		case IsWithdrawalPossible(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.WithdrawalPossible
		case IsExitedSlashed(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.ExitedSlashed
		case IsExitedUnslashed(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.ExitedUnslashed
		case IsActiveSlashed(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.ActiveSlashed
		case IsActiveExited(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.ActiveExited
		case IsActiveOngoing(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.ActiveOngoing
		case IsPendingQueued(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.PendingQueued
		case IsPendingInitialized(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.PendingInitialized
		case IsRegisteredToStakewise(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.RegisteredToStakewise
		case IsUploadedStakewise(pubKey, statuses):
			validatorStatuses[pubKey] = swapi.UploadedStakewise
		case IsUploadedToNodeset(pubKey, statuses, registeredPubkeys):
			validatorStatuses[pubKey] = swapi.UploadedToNodeset
		default:
			validatorStatuses[pubKey] = swapi.Generated
		}
	}

	data.ValidatorStatus = validatorStatuses

	return nil
}

func IsPendingInitialized(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_PendingInitialized
}
func IsPendingQueued(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_PendingQueued
}
func IsActiveOngoing(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ActiveOngoing
}

func IsActiveExited(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ActiveExiting
}

func IsActiveSlashed(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ActiveSlashed
}

func IsExitedUnslashed(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ExitedUnslashed
}

func IsExitedSlashed(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ExitedSlashed
}

func IsWithdrawalPossible(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_WithdrawalPossible
}

func IsWithdrawalDone(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_WithdrawalDone
}

func IsRegisteredToStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

func IsUploadedStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

func IsUploadedToNodeset(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus, registeredPubkeys []beacon.ValidatorPubkey) bool {
	for _, registeredPubKey := range registeredPubkeys {
		if registeredPubKey == pubKey {
			return true
		}
	}
	return false
}
