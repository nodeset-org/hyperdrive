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

type stakeWiseGetDepositDataSetVersionContextFactory struct {
	handler *StakeWiseHandler
}

func (f *stakeWiseGetDepositDataSetVersionContextFactory) Create(args url.Values) (*stakeWiseGetDepositDataSetVersionContext, error) {
	c := &stakeWiseGetDepositDataSetVersionContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("vault", args, input.ValidateAddress, &c.vault),
	}
	return c, errors.Join(inputErrs...)
}

func (f *stakeWiseGetDepositDataSetVersionContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*stakeWiseGetDepositDataSetVersionContext, api.NodeSetStakeWise_GetDepositDataSetVersionData](
		router, "get-deposit-data-set/version", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseGetDepositDataSetVersionContext struct {
	handler *StakeWiseHandler
	vault   common.Address
}

func (c *stakeWiseGetDepositDataSetVersionContext) PrepareData(data *api.NodeSetStakeWise_GetDepositDataSetVersionData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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

	// Get the set version
	ns := sp.GetNodeSetServiceManager()
	version, err := ns.StakeWise_GetServerDepositDataVersion(ctx, c.vault)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	data.Version = version
	return types.ResponseStatus_Success, nil
}
