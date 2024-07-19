package ns_constellation

import (
	"errors"
	"net/url"

	hdcommon "github.com/nodeset-org/hyperdrive-daemon/common"
	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	apiv2 "github.com/nodeset-org/nodeset-client-go/api-v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/utils/input"
)

// ===============
// === Factory ===
// ===============

type constellationGetDepositSignatureContextFactory struct {
	handler *ConstellationHandler
}

func (f *constellationGetDepositSignatureContextFactory) Create(args url.Values) (*constellationGetDepositSignatureContext, error) {
	c := &constellationGetDepositSignatureContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("minipoolAddress", args, input.ValidateAddress, &c.minipoolAddress),
		server.ValidateArg("salt", args, input.ValidateByteArray, &c.salt),
	}
	return c, errors.Join(inputErrs...)
}

func (f *constellationGetDepositSignatureContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*constellationGetDepositSignatureContext, api.NodeSetConstellation_GetDepositSignatureData](
		router, "get-deposit-signature", f, f.handler.logger.Logger, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============
type constellationGetDepositSignatureContext struct {
	handler         *ConstellationHandler
	minipoolAddress common.Address
	salt            []byte
}

func (c *constellationGetDepositSignatureContext) PrepareData(data *api.NodeSetConstellation_GetDepositSignatureData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
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
	timestamp, signature, err := ns.Constellation_GetDepositSignatureAndTime(ctx, c.minipoolAddress, c.salt)
	if err != nil {
		if errors.Is(err, apiv2.ErrNotAuthorized) {
			data.NotAuthorized = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, apiv2.ErrMinipoolLimitReached) {
			data.LimitReached = true
			return types.ResponseStatus_Success, nil
		}
		if errors.Is(err, apiv2.ErrMissingExitMessage) {
			data.MissingExitMessage = true
			return types.ResponseStatus_Success, nil
		}
		return types.ResponseStatus_Error, err
	}

	data.Signature = signature
	data.Time = timestamp
	return types.ResponseStatus_Success, nil
}
