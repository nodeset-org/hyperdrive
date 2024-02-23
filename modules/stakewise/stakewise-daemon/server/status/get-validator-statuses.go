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
	privateKeys, err := w.GetAllPrivateKeys()
	if err != nil {
		return fmt.Errorf("error getting private keys: %w", err)
	}
	publicKeys, err := w.DerivePubKeys(privateKeys)
	if err != nil {
		return fmt.Errorf("error getting public keys: %w", err)
	}
	statuses, err := bc.GetValidatorStatuses(context.Background(), publicKeys, nil)
	fmt.Printf("!!! statuses: %v\n", statuses)
	if err != nil {
		return fmt.Errorf("error getting validator statuses: %w", err)
	}
	var generatedValidators, uploadedNodesetValidators, uploadedStakewiseValidators, registeredStakewiseValidators, withdrawalDone, withdrawalPossible, exitedSlashed, exitedUnslashed, activeSlashed, activeExited, activeOngoing, pendingQueued, pendingInitialized []beacon.ValidatorPubkey

	for _, pubKey := range publicKeys {
		if IsWithdrawalDone(pubKey, statuses) {
			withdrawalDone = append(withdrawalDone, pubKey)
		} else if IsWithdrawalPossible(pubKey, statuses) {
			withdrawalPossible = append(withdrawalPossible, pubKey)
		} else if IsExitedSlashed(pubKey, statuses) {
			exitedSlashed = append(exitedSlashed, pubKey)
		} else if IsExitedUnslashed(pubKey, statuses) {
			exitedUnslashed = append(exitedUnslashed, pubKey)
		} else if IsActiveSlashed(pubKey, statuses) {
			activeSlashed = append(activeSlashed, pubKey)
		} else if IsActiveExited(pubKey, statuses) {
			activeExited = append(activeExited, pubKey)
		} else if IsActiveOngoing(pubKey, statuses) {
			activeOngoing = append(activeOngoing, pubKey)
		} else if IsPendingQueued(pubKey, statuses) {
			pendingQueued = append(pendingQueued, pubKey)
		} else if IsPendingInitialized(pubKey, statuses) {
			pendingInitialized = append(pendingInitialized, pubKey)
		} else if IsRegisteredToStakewise(pubKey, statuses) {
			registeredStakewiseValidators = append(registeredStakewiseValidators, pubKey)
		} else if IsUploadedStakewise(pubKey, statuses) {
			uploadedStakewiseValidators = append(uploadedStakewiseValidators, pubKey)
		} else if IsUploadedToNodeset(pubKey, statuses) {
			uploadedNodesetValidators = append(uploadedNodesetValidators, pubKey)
		} else if IsGenerated(pubKey, statuses) {
			generatedValidators = append(generatedValidators, pubKey)
		} else {
			fmt.Printf("Unknown status for validator %s\n", pubKey.HexWithPrefix())
		}
	}

	data.Generated = generatedValidators
	data.UploadedToNodeset = uploadedNodesetValidators
	data.UploadToStakewise = uploadedStakewiseValidators
	data.RegisteredToStakewise = registeredStakewiseValidators
	data.PendingInitialized = pendingInitialized
	data.PendingQueued = pendingQueued
	data.ActiveOngoing = activeOngoing
	data.ActiveExited = activeExited
	data.ActiveSlashed = activeSlashed
	data.ExitedUnslashed = exitedUnslashed
	data.ExitedSlashed = exitedSlashed
	data.WithdrawalPossible = withdrawalPossible
	data.WithdrawalDone = withdrawalDone

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

// IMPORTANT
func IsUploadedToNodeset(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	// TODO: Implement
	return false
}

// IMPORTANT
func IsGenerated(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return true
}
