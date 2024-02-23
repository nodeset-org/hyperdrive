package swstatus

import (
	"errors"
	"fmt"
	"net/url"

	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
)

// ===============
// === Factory ===
// ===============

type statusGetActiveValidatorsContextFactory struct {
	handler *StatusHandler
}

func (f *statusGetActiveValidatorsContextFactory) Create(args url.Values) (*statusGetActiveValidatorsContext, error) {
	c := &statusGetActiveValidatorsContext{
		handler: f.handler,
	}
	inputErrs := []error{}
	return c, errors.Join(inputErrs...)
}

func (f *statusGetActiveValidatorsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*statusGetActiveValidatorsContext, swapi.ActiveValidatorsData](
		router, "status", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type statusGetActiveValidatorsContext struct {
	handler *StatusHandler
}

func (c *statusGetActiveValidatorsContext) PrepareData(data *swapi.ActiveValidatorsData, opts *bind.TransactOpts) error {
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
	data.ActiveValidators = publicKeys
	return nil
}
