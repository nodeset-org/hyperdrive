package nodeset

import (
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type nodeSetRegisterNodeContextFactory struct {
	handler *NodeSetHandler
}

func (f *nodeSetRegisterNodeContextFactory) Create(args url.Values) (*nodeSetRegisterNodeContext, error) {
	c := &nodeSetRegisterNodeContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.GetStringFromVars("email", args, &c.email),
	}
	return c, errors.Join(inputErrs...)
}

func (f *nodeSetRegisterNodeContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*nodeSetRegisterNodeContext, api.NodeSetRegisterNodeData](
		router, "register-node", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodeSetRegisterNodeContext struct {
	handler *NodeSetHandler
	email   string
}

func (c *nodeSetRegisterNodeContext) PrepareData(data *api.NodeSetRegisterNodeData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ctx := c.handler.ctx

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}

	// Register the node
	ns := sp.GetNodeSetServiceManager()
	result, err := ns.RegisterNode(ctx, c.email)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	// Handle the result options
	switch result {
	case common.RegistrationResult_Success:
		data.Success = true
		// Force a re-login to update the registration status
		_ = sp.GetNodeSetServiceManager().Login(ctx)
	case common.RegistrationResult_AlreadyRegistered:
		data.Success = false
		data.AlreadyRegistered = true
	case common.RegistrationResult_NotWhitelisted:
		data.Success = false
		data.NotWhitelisted = true
	}
	return types.ResponseStatus_Success, nil
}
