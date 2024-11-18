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

type stakeWiseGetDepositDataSetContextFactory struct {
	handler *StakeWiseHandler
}

func (f *stakeWiseGetDepositDataSetContextFactory) Create(args url.Values) (*stakeWiseGetDepositDataSetContext, error) {
	c := &stakeWiseGetDepositDataSetContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.GetStringFromVars("deployment", args, &c.deployment),
		server.ValidateArg("vault", args, input.ValidateAddress, &c.vault),
	}
	return c, errors.Join(inputErrs...)
}

func (f *stakeWiseGetDepositDataSetContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*stakeWiseGetDepositDataSetContext, api.NodeSetStakeWise_GetDepositDataSetData](
		router, "get-deposit-data-set", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============
type stakeWiseGetDepositDataSetContext struct {
	handler *StakeWiseHandler

	deployment string
	vault      common.Address
}

func (c *stakeWiseGetDepositDataSetContext) PrepareData(data *api.NodeSetStakeWise_GetDepositDataSetData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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
	version, set, err := ns.StakeWise_GetServerDepositData(ctx, c.deployment, c.vault)
	if err != nil {
		return types.ResponseStatus_Error, err
	}

	data.Version = version
	data.DepositData = set
	return types.ResponseStatus_Success, nil
}
