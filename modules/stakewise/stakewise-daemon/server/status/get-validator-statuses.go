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
	var generatedValidators, uploadedNodesetValidators, uploadedStakewiseValidators, registeredStakewiseValidators, waitingDepositConfirmationValidators, depositingValidators, depositedValidators, activeValidators, exitingValidators, exitedValidators []beacon.ValidatorPubkey

	for _, pubKey := range publicKeys {
		if IsExited(pubKey, statuses) {
			exitedValidators = append(exitedValidators, pubKey)
		} else if IsExiting(pubKey, statuses) {
			exitingValidators = append(exitingValidators, pubKey)
		} else if IsActive(pubKey, statuses) {
			activeValidators = append(activeValidators, pubKey)
		} else if IsDeposited(pubKey, statuses) {
			depositedValidators = append(depositedValidators, pubKey)
		} else if IsDepositing(pubKey, statuses) {
			depositingValidators = append(depositingValidators, pubKey)
		} else if IsWaitingDepositConfirmation(pubKey, statuses) {
			waitingDepositConfirmationValidators = append(waitingDepositConfirmationValidators, pubKey)
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
	data.WaitingDepositConfirmation = waitingDepositConfirmationValidators
	data.Depositing = depositingValidators
	data.Deposited = depositedValidators
	data.Active = activeValidators
	data.Exiting = exitingValidators
	data.Exited = exitedValidators

	return nil
}

func IsExited(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ExitedSlashed || statuses[pubKey].Status == types.ValidatorState_ExitedUnslashed
}

func IsExiting(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ActiveExiting || statuses[pubKey].Status == types.ValidatorState_WithdrawalPossible || statuses[pubKey].Status == types.ValidatorState_WithdrawalDone
}

func IsActive(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_ActiveOngoing
}

func IsDeposited(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return false
}

func IsDepositing(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_PendingInitialized
}

func IsWaitingDepositConfirmation(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return statuses[pubKey].Status == types.ValidatorState_PendingQueued
}

func IsRegisteredToStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return false
}

func IsUploadedStakewise(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return false
}

func IsUploadedToNodeset(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return false
}

func IsGenerated(pubKey beacon.ValidatorPubkey, statuses map[beacon.ValidatorPubkey]types.ValidatorStatus) bool {
	return true
}
