package ns_stakewise

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/rocket-pool/node-manager-core/utils/input"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type stakeWiseGetRegisteredValidatorsContextFactory struct {
	handler *StakeWiseHandler
}

func (f *stakeWiseGetRegisteredValidatorsContextFactory) Create(args url.Values) (*stakeWiseGetRegisteredValidatorsContext, error) {
	c := &stakeWiseGetRegisteredValidatorsContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("vault", args, input.ValidateAddress, &c.vault),
	}
	return c, errors.Join(inputErrs...)
}

func (f *stakeWiseGetRegisteredValidatorsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*stakeWiseGetRegisteredValidatorsContext, api.NodeSetStakeWise_GetRegisteredValidatorsData](
		router, "get-registered-validators", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseGetRegisteredValidatorsContext struct {
	handler *StakeWiseHandler
	vault   common.Address
}

func (c *stakeWiseGetRegisteredValidatorsContext) PrepareData(data *api.NodeSetStakeWise_GetRegisteredValidatorsData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ctx := c.handler.ctx

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}
	err = sp.RequireRegisteredWithNodeSet(ctx)
	if err != nil {
		if errors.Is(err, hdcommon.ErrNotRegisteredWithNodeSet) {
			data.NotRegistered = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	// Get the registered validators
	ns := sp.GetNodeSetServiceManager()
	response, err := ns.StakeWise_GetRegisteredValidators(ctx, c.vault)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	data.Validators = response
	return types.ResponseStatus_Success, nil
}
