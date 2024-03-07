package utils

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_input "github.com/rocket-pool/node-manager-core/utils/input"
	ens "github.com/wealdtech/go-ens/v3"
)

// ===============
// === Factory ===
// ===============

type utilsResolveEnsContextFactory struct {
	handler *UtilsHandler
}

func (f *utilsResolveEnsContextFactory) Create(args url.Values) (*utilsResolveEnsContext, error) {
	c := &utilsResolveEnsContext{
		handler: f.handler,
	}
	inputErrs := []error{
		nmc_server.ValidateArg("address", args, nmc_input.ValidateAddress, &c.address),
		nmc_server.GetStringFromVars("name", args, &c.name),
	}
	return c, errors.Join(inputErrs...)
}

func (f *utilsResolveEnsContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*utilsResolveEnsContext, api.UtilsResolveEnsData](
		router, "resolve-ens", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type utilsResolveEnsContext struct {
	handler *UtilsHandler
	address common.Address
	name    string
}

func (c *utilsResolveEnsContext) PrepareData(data *api.UtilsResolveEnsData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()

	emptyAddress := common.Address{}
	if c.address != emptyAddress {
		data.Address = c.address
		name, err := ens.ReverseResolve(ec, c.address)
		if err != nil {
			data.FormattedName = data.Address.Hex()
		} else {
			data.EnsName = name
			data.FormattedName = fmt.Sprintf("%s (%s)", name, data.Address.Hex())
		}
	} else if c.name != "" {
		data.EnsName = c.name
		address, err := ens.Resolve(ec, c.name)
		if err != nil {
			return fmt.Errorf("error resolving ENS address for [%s]: %w", c.name, err)
		}
		data.Address = address
		data.FormattedName = fmt.Sprintf("%s (%s)", c.name, data.Address.Hex())
	} else {
		return fmt.Errorf("either address or name must not be empty")
	}

	return nil
}
