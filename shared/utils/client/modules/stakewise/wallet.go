package swclient

import (
	"strconv"

	"github.com/nodeset-org/hyperdrive/shared/types/api"
	swapi "github.com/nodeset-org/hyperdrive/shared/types/api/modules/stakewise"
	"github.com/nodeset-org/hyperdrive/shared/utils/client"
)

type WalletRequester struct {
	context *client.RequesterContext
}

func NewWalletRequester(context *client.RequesterContext) *WalletRequester {
	return &WalletRequester{
		context: context,
	}
}

func (r *WalletRequester) GetName() string {
	return "Wallet"
}
func (r *WalletRequester) GetRoute() string {
	return "wallet"
}
func (r *WalletRequester) GetContext() *client.RequesterContext {
	return r.context
}

// Generate and save new validator keys
func (r *WalletRequester) GenerateKeys(count uint64) (*api.ApiResponse[swapi.WalletGenerateKeysData], error) {
	args := map[string]string{
		"count": strconv.FormatUint(count, 10),
	}
	return client.SendGetRequest[swapi.WalletGenerateKeysData](r, "generate-keys", "GenerateKeys", args)
}

// Export the wallet in encrypted ETH key format
func (r *WalletRequester) Initialize() (*api.ApiResponse[swapi.WalletInitializeData], error) {
	return client.SendGetRequest[swapi.WalletInitializeData](r, "initialize", "Initialize", nil)
}

// Regenerate the aggregated deposit data file for Stakewise to use
func (r *WalletRequester) RegenerateDepositData() (*api.ApiResponse[swapi.WalletRegenerateDepositDataData], error) {
	return client.SendGetRequest[swapi.WalletRegenerateDepositDataData](r, "regen-deposit-data", "RegenerateDepositData", nil)
}
