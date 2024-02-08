package swstatus

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
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
	inputErrs := []error{
		// server.ValidateArg("root", args, input.ValidateHash, &c.root),
	}
	return c, errors.Join(inputErrs...)
}

func (f *statusGetActiveValidatorsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*statusGetActiveValidatorsContext, api.ActiveValidatorsData](
		router, "status", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type statusGetActiveValidatorsContext struct {
	handler *StatusHandler
	// root    common.Hash
}

func (c *statusGetActiveValidatorsContext) PrepareData(data *api.ActiveValidatorsData, opts *bind.TransactOpts) error {
	fmt.Printf("statusGetActiveValidatorsContext.PrepareData data: %+v\n", data)
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
	fmt.Printf("statusGetActiveValidatorsContext.PrepareData publicKeys: %+v\n", publicKeys)
	return nil
}
