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

type walletSignMessageContextFactory struct {
	handler *WalletHandler
}

func (f *walletSignMessageContextFactory) Create(args url.Values) (*walletSignMessageContext, error) {
	c := &walletSignMessageContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("message", args, input.ValidateByteArray, &c.message),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSignMessageContextFactory) RegisterRoute(router *mux.Router) {
	utils.RegisterQuerylessGet[*walletSignMessageContext, api.WalletSignMessageData](
		router, "sign-message", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletSignMessageContext struct {
	handler *WalletHandler
	message []byte
}

func (c *walletSignMessageContext) PrepareData(data *api.WalletSignMessageData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	err := errors.Join(
		sp.RequireWalletReady(),
	)
	if err != nil {
		return err
	}

	signedBytes, err := w.SignMessage(c.message)
	if err != nil {
		return fmt.Errorf("error signing message: %w", err)
	}
	data.SignedMessage = signedBytes
	return nil
}
