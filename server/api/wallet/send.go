package wallet

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/eth/contracts"

	"github.com/nodeset-org/hyperdrive-daemon/shared/types/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/utils/input"
)

// ===============
// === Factory ===
// ===============

type walletSendContextFactory struct {
	handler *WalletHandler
}

func (f *walletSendContextFactory) Create(args url.Values) (*walletSendContext, error) {
	c := &walletSendContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("amount", args, input.ValidateBigInt, &c.amount),
		server.GetStringFromVars("token", args, &c.token),
		server.ValidateArg("recipient", args, input.ValidateAddress, &c.recipient),
	}
	return c, errors.Join(inputErrs...)
}

func (f *walletSendContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletSendContext, api.WalletSendData](
		router, "send", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletSendContext struct {
	handler *WalletHandler

	amount    *big.Int
	token     string
	recipient common.Address
}

func (c *walletSendContext) PrepareData(data *api.WalletSendData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ec := sp.GetEthClient()
	qMgr := sp.GetQueryManager()
	txMgr := sp.GetTransactionManager()
	ctx := c.handler.ctx
	nodeAddress, _ := sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireNodeAddress()
	if err != nil {
		return types.ResponseStatus_AddressNotPresent, err
	}
	err = sp.RequireEthClientSynced(ctx)
	if err != nil {
		return types.ResponseStatus_AddressNotPresent, err
	}

	// Get the contract (nil in the case of ETH)
	var tokenContract contracts.IErc20Token
	if c.token == "eth" {
		tokenContract = nil
	} else if strings.HasPrefix(c.token, "0x") {
		// Arbitrary token - make sure the contract address is legal
		if !common.IsHexAddress(c.token) {
			return types.ResponseStatus_InvalidArguments, fmt.Errorf("[%s] is not a valid token address", c.token)
		}
		tokenAddress := common.HexToAddress(c.token)

		// Make a binding for it
		tokenContract, err := contracts.NewErc20Contract(tokenAddress, ec, qMgr, txMgr, nil)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error creating ERC20 contract binding: %w", err)
		}
		data.TokenSymbol = tokenContract.Symbol()
		data.TokenName = tokenContract.Name()
	}

	// Get the balance
	if tokenContract != nil {
		err := qMgr.Query(func(mc *batch.MultiCaller) error {
			tokenContract.BalanceOf(mc, &data.Balance, nodeAddress)
			return nil
		}, nil)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting token balance: %w", err)
		}
	} else {
		// ETH balance
		var err error
		data.Balance, err = ec.BalanceAt(context.Background(), nodeAddress, nil)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting ETH balance: %w", err)
		}
	}

	// Check the balance
	data.InsufficientBalance = (data.Balance.Cmp(common.Big0) == 0)
	data.CanSend = !(data.InsufficientBalance)

	// Get the TX Info
	if data.CanSend {
		var txInfo *eth.TransactionInfo
		var err error
		if tokenContract != nil {
			txInfo, err = tokenContract.Transfer(c.recipient, c.amount, opts)
		} else {
			// ETH transfers
			newOpts := &bind.TransactOpts{
				From:  opts.From,
				Nonce: opts.Nonce,
				Value: c.amount,
			}
			txInfo = txMgr.CreateTransactionInfoRaw(c.recipient, nil, newOpts)
		}
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting TX info for Transfer: %w", err)
		}
		data.TxInfo = txInfo
	}

	return types.ResponseStatus_Success, nil
}
