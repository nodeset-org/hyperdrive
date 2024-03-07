package wallet

import (
	"errors"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	nmc_server "github.com/rocket-pool/node-manager-core/api/server"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
	nmc_input "github.com/rocket-pool/node-manager-core/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletSendMessageContextFactory struct {
	handler *WalletHandler
}

func (f *walletSendMessageContextFactory) Create(args url.Values) (*walletSendMessageContext, error) {
	c := &walletSendMessageContext{
		handler: f.handler,
	}
	inputErrs := []error{
		nmc_server.ValidateArg("message", args, nmc_input.ValidateByteArray, &c.message),
		nmc_server.ValidateArg("address", args, nmc_input.ValidateAddress, &c.address),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSendMessageContextFactory) RegisterRoute(router *mux.Router) {
	nmc_server.RegisterQuerylessGet[*walletSendMessageContext, nmc_types.TxInfoData](
		router, "send-message", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletSendMessageContext struct {
	handler *WalletHandler
	message []byte
	address common.Address
}

func (c *walletSendMessageContext) PrepareData(data *nmc_types.TxInfoData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	txMgr := sp.GetTransactionManager()

	err := errors.Join(
		sp.RequireWalletReady(),
	)
	if err != nil {
		return err
	}

	txInfo := txMgr.CreateTransactionInfoRaw(c.address, c.message, opts)
	data.TxInfo = txInfo
	return nil
}
