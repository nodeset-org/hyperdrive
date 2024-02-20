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

	var activeValidators, exitingValidators, exitedValidators []beacon.ValidatorPubkey

	for _, pubKey := range publicKeys {
		if IsExiting(pubKey) { // Assume IsExiting is a function you define to check if a validator is exiting
			exitingValidators = append(exitingValidators, pubKey)
		} else if IsExited(pubKey) { // Assume IsExited is a function you define to check if a validator has exited
			exitedValidators = append(exitedValidators, pubKey)
		} else {
			activeValidators = append(activeValidators, pubKey)
		}
	}

	data.ActiveValidators = activeValidators
	data.ExitingValidators = exitingValidators
	data.ExitedValidators = exitedValidators
	return nil
}

func IsExiting(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}

func IsExited(pubKey beacon.ValidatorPubkey) bool {
	// TODO
	return true
}
