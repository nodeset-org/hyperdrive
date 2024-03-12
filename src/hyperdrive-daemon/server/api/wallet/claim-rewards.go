package wallet

import (
	"fmt"
	"net/url"
	"strings"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	localABI "github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/service/abi"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/utils"
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
	utils.RegisterQuerylessGet[*walletClaimRewardsContext, api.SuccessData](
		router, "claim-rewards", f, f.handler.serviceProvider,
	)
}

const YieldDistributorContractAddress = "0x"

// ===============
// === Context ===
// ===============

type walletClaimRewardsContext struct {
	handler *WalletHandler
	// address common.Address
}

func (c *walletClaimRewardsContext) PrepareData(data *api.SuccessData, opts *bind.TransactOpts) error {
	fmt.Printf("Preparing data for claim reward\n")
	sp := c.handler.serviceProvider
	// w := sp.GetWallet()
	// TODO: HUY!!!
	ec := sp.GetEthClient()

	abi, err := abi.JSON(strings.NewReader(localABI.YieldDistributorABI))
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(YieldDistributorContractAddress)
	boundContract := bind.NewBoundContract(contractAddress, abi, ec, ec, ec)
	fmt.Printf("Bound contract: %v\n", boundContract)
	// ec.SendTransaction(opts, w.Address)

	return nil
}
