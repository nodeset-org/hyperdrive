package wallet

import (
	"errors"
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletSignTxContextFactory struct {
	handler *WalletHandler
}

func (f *walletSignTxContextFactory) Create(args url.Values) (*walletSignTxContext, error) {
	c := &walletSignTxContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("tx", args, input.ValidateByteArray, &c.tx),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSignTxContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletSignTxContext, api.WalletSignTxData](
		router, "sign-tx", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletSignTxContext struct {
	handler *WalletHandler
	tx      []byte
}

func (c *walletSignTxContext) PrepareData(data *api.WalletSignTxData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	err := errors.Join(
		sp.RequireWalletReady(),
	)
	if err != nil {
		return err
	}

	signedBytes, err := w.SignTransaction(c.tx)
	if err != nil {
		return fmt.Errorf("error signing transaction: %w", err)
	}
	data.SignedTx = signedBytes
	return nil
}
