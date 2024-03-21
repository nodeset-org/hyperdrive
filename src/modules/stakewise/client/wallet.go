package swclient

import (
	"strconv"

	"github.com/nodeset-org/hyperdrive/client"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
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
func (r *WalletRequester) GenerateKeys(count uint64, restartVc bool) (*api.ApiResponse[swapi.WalletGenerateKeysData], error) {
	args := map[string]string{
		"count":      strconv.FormatUint(count, 10),
		"restart-vc": strconv.FormatBool(restartVc),
	}
	return client.SendGetRequest[swapi.WalletGenerateKeysData](r, "generate-keys", "GenerateKeys", args)
}

func (r *WalletRequester) ClaimRewards() (*api.ApiResponse[api.TxInfoData], error) {
	return client.SendGetRequest[api.TxInfoData](r, "claim-rewards", "ClaimRewards", nil)
}

// Export the wallet in encrypted ETH key format
func (r *WalletRequester) Initialize() (*api.ApiResponse[swapi.WalletInitializeData], error) {
	return client.SendGetRequest[swapi.WalletInitializeData](r, "initialize", "Initialize", nil)
}
