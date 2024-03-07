package swclient

import (
	"strconv"

	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	nmc_client "github.com/rocket-pool/node-manager-core/api/client"
	nmc_types "github.com/rocket-pool/node-manager-core/api/types"
)

type WalletRequester struct {
	context *nmc_client.RequesterContext
}

func NewWalletRequester(context *nmc_client.RequesterContext) *WalletRequester {
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
func (r *WalletRequester) GetContext() *nmc_client.RequesterContext {
	return r.context
}

// Generate and save new validator keys
func (r *WalletRequester) GenerateKeys(count uint64, restartVc bool) (*nmc_types.ApiResponse[swapi.WalletGenerateKeysData], error) {
	args := map[string]string{
		"count":      strconv.FormatUint(count, 10),
		"restart-vc": strconv.FormatBool(restartVc),
	}
	return nmc_client.SendGetRequest[swapi.WalletGenerateKeysData](r, "generate-keys", "GenerateKeys", args)
}

// Export the wallet in encrypted ETH key format
func (r *WalletRequester) Initialize() (*nmc_types.ApiResponse[swapi.WalletInitializeData], error) {
	return nmc_client.SendGetRequest[swapi.WalletInitializeData](r, "initialize", "Initialize", nil)
}
