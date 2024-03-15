package wallet

import (
	"fmt"
	"math/big"
	"net/url"
	"strings"
	_ "time/tzdata"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/server/api/contract"
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

const SplitMainAddress = "0x2ed6c4B5dA6378c7897AC67Ba9e43102Feb694EE"

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
	w := sp.GetWallet()
	walletAddress, _ := w.GetAddress()
	// TODO: HUY!!!
	ec := sp.GetEthClient()

	abi, err := abi.JSON(strings.NewReader(localABI.SplitMainABI))
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(SplitMainAddress)
	boundContract := bind.NewBoundContract(contractAddress, abi, ec, ec, ec)
	contractInstance := &contract.Contract{
		Contract: boundContract,
		Address:  &contractAddress,
		ABI:      &abi,
		Client:   ec,
	}
	opts, err = w.GetTransactor()
	if err != nil {
		return err
	}
	fmt.Printf("Contract instance: %v\n", contractInstance)
	tx, err := contractInstance.Transact(opts, "withdraw", walletAddress, big.NewInt(0), []common.Address{})
	if err != nil {
		return err
	}
	fmt.Printf("Transaction: %s\n", tx)
	return nil
}
