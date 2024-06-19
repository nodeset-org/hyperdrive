package wallet

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

type walletExportContextFactory struct {
	handler *WalletHandler
}

func (f *walletExportContextFactory) Create(args url.Values) (*walletExportContext, error) {
	c := &walletExportContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletExportContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletExportContext, api.WalletExportData](
		router, "export", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletExportContext struct {
	handler *WalletHandler
}

func (c *walletExportContext) PrepareData(data *api.WalletExportData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return types.ResponseStatus_WalletNotReady, err
	}

	// Get password - blank if not saved
	pw, isSet, err := w.GetPassword()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting wallet password: %w", err)
	}
	if isSet {
		data.Password = pw
	}

	// Serialize wallet
	walletString, err := w.SerializeData()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error serializing wallet keystore: %w", err)
	}
	data.Wallet = walletString

	// Get account private key
	data.AccountPrivateKey, err = w.GetNodePrivateKeyBytes()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting node wallet private key: %w", err)
	}
	return types.ResponseStatus_Success, nil
}
