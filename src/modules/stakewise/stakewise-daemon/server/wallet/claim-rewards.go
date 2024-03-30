package swwallet

import (
	"fmt"
	"net/url"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swcontracts "github.com/nodeset-org/hyperdrive/modules/stakewise/stakewise-daemon/common/contracts"
	"github.com/rocket-pool/node-manager-core/api/types"
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

	return c, nil
}

func (f *walletClaimRewardsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletClaimRewardsContext, types.TxInfoData](
		router, "claim-rewards", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
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
func (c *walletClaimRewardsContext) PrepareData(data *types.TxInfoData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	logger := c.handler.logger
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	res := sp.GetResources()
	txMgr := sp.GetTransactionManager()

	logger.Debug("Preparing data for claim reward")
	splitMainContract, err := swcontracts.NewSplitMain(res.Splitmain, ec, txMgr)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error creating Stakewise Vault binding: %w", err)
	}

	data.TxInfo, err = splitMainContract.Withdraw(res.Vault, *res.ClaimEthAmount, res.ClaimTokenList, opts)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error creating Withdraw TX: %w", err)
	}
	return types.ResponseStatus_Success, nil
}
