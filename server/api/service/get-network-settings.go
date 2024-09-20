package service

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type serviceGetNetworkSettingsContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceGetNetworkSettingsContextFactory) Create(args url.Values) (*serviceGetNetworkSettingsContext, error) {
	c := &serviceGetNetworkSettingsContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *serviceGetNetworkSettingsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*serviceGetNetworkSettingsContext, api.ServiceGetNetworkSettingsData](
		router, "get-network-settings", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceGetNetworkSettingsContext struct {
	handler *ServiceHandler
}

func (c *serviceGetNetworkSettingsContext) PrepareData(data *api.ServiceGetNetworkSettingsData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	cfg := sp.GetConfig()
	settingsList := cfg.GetNetworkSettings()
	network := cfg.Network.Value
	for _, settings := range settingsList {
		if settings.Key == network {
			data.Settings = settings
			return types.ResponseStatus_Success, nil
		}
	}
	return types.ResponseStatus_Error, fmt.Errorf("hyperdrive has network [%s] selected but there are no settings for it", network)
}
