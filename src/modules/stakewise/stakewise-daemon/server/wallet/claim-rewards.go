package swwallet

import (
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcontracts "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common/contracts"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type walletClaimRewardsContextFactory struct {
	handler *WalletHandler
}

func (f *walletClaimRewardsContextFactory) Create(args url.Values) (*walletClaimRewardsContext, error) {
	c := &walletClaimRewardsContext{
		handler: f.handler,
	}
	// inputErrs := []error{
	// 	server.ValidateArg("address", args, input.ValidateAddress, &c.address),
	// }
	return c, nil
}

func (f *walletClaimRewardsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletClaimRewardsContext, api.TxInfoData](
		router, "claim-rewards", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletClaimRewardsContext struct {
	handler *WalletHandler
	// address common.Address
}

// Return the transaction data
func (c *walletClaimRewardsContext) PrepareData(data *api.TxInfoData, opts *bind.TransactOpts) error {
	fmt.Printf("Preparing data for claim reward\n")

	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	res := sp.GetResources()
	txMgr := sp.GetTransactionManager()

	splitMainContract, err := swcontracts.NewSplitMain(res.Splitmain, ec, txMgr)
	if err != nil {
		return fmt.Errorf("error creating Stakewise Vault binding: %w", err)
	}

	data.TxInfo, err = splitMainContract.SetWithdraw(res.Vault, opts)
	if err != nil {
		return fmt.Errorf("error creating SetWithdraw TX: %w", err)
	}
	return nil
}
