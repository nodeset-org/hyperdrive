package swstatus

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/nodeset-org/eth-utils/beacon"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"

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
	w := sp.GetWallet()
	privateKeys, err := w.GetAllPrivateKeys()
	if err != nil {
		return fmt.Errorf("error getting private keys: %w", err)
	}
	publicKeys, err := w.DerivePubKeys(privateKeys)
	if err != nil {
		return fmt.Errorf("error getting public keys: %w", err)
	}

	var generatedValidators, uploadedNodesetValidators, uploadedStakewiseValidators, registeredStakewiseValidators, waitingDepositConfirmationValidators, depositingValidators, depositedValidators, activeValidators, exitingValidators, exitedValidators []beacon.ValidatorPubkey

	for _, pubKey := range publicKeys {
		if IsExited(pubKey) {
			exitedValidators = append(exitedValidators, pubKey)
		} else if IsExiting(pubKey) {
			exitingValidators = append(exitingValidators, pubKey)
		} else if IsActive(pubKey) {
			activeValidators = append(activeValidators, pubKey)
		} else if IsDeposited(pubKey) {
			depositedValidators = append(depositedValidators, pubKey)
		} else if IsDepositing(pubKey) {
			depositingValidators = append(depositingValidators, pubKey)
		} else if IsWaitingDepositConfirmation(pubKey) {
			waitingDepositConfirmationValidators = append(waitingDepositConfirmationValidators, pubKey)
		} else if IsRegisteredToStakewise(pubKey) {
			registeredStakewiseValidators = append(registeredStakewiseValidators, pubKey)
		} else if IsUploadedStakewise(pubKey) {
			uploadedStakewiseValidators = append(uploadedStakewiseValidators, pubKey)
		} else if IsUploadedToNodeset(pubKey) {
			uploadedNodesetValidators = append(uploadedNodesetValidators, pubKey)
		} else if IsGenerated(pubKey) {
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

func IsExited(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsExiting(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsActive(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsDeposited(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsDepositing(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsWaitingDepositConfirmation(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsRegisteredToStakewise(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsUploadedStakewise(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsUploadedToNodeset(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsGenerated(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}
