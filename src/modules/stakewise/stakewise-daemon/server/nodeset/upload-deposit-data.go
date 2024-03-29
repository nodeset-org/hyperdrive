package swnodeset

import (
	"errors"
	"fmt"
	"math/big"
	"net/url"

	"github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	duserver "github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/node-manager-core/utils/input"
	eth2types "github.com/wealdtech/go-eth2-types/v2"
)

// ===============
// === Factory ===
// ===============

type nodesetUploadDepositDataContextFactory struct {
	handler *NodesetHandler
}

func (f *nodesetUploadDepositDataContextFactory) Create(args url.Values) (*nodesetUploadDepositDataContext, error) {
	c := &nodesetUploadDepositDataContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("bypassBalanceCheck", args, input.ValidateBool, &c.bypassBalanceCheck),
	}
	return c, errors.Join(inputErrs...)
}

func (f *nodesetUploadDepositDataContextFactory) RegisterRoute(router *mux.Router) {
	duserver.RegisterQuerylessGet[*nodesetUploadDepositDataContext, swapi.NodesetUploadDepositDataData](
		router, "upload-deposit-data", f, f.handler.logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodesetUploadDepositDataContext struct {
	handler            *NodesetHandler
	bypassBalanceCheck bool
}

func (c *nodesetUploadDepositDataContext) PrepareData(data *swapi.NodesetUploadDepositDataData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	ddMgr := sp.GetDepositDataManager()
	nc := sp.GetNodesetClient()
	w := sp.GetWallet()
	ec := sp.GetEthClient()
	ctx := c.handler.ctx

	// Get the list of registered validators
	registeredPubkeyMap := map[beacon.ValidatorPubkey]bool{}
	pubkeyStatusResponse, err := nc.GetRegisteredValidators()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting registered validators: %w", err)
	}
	registeredPubkeys := []beacon.ValidatorPubkey{}
	for _, pubkeyStatus := range pubkeyStatusResponse {
		registeredPubkeys = append(registeredPubkeys, pubkeyStatus.Pubkey)
	}
	for _, pubkey := range registeredPubkeys {
		registeredPubkeyMap[pubkey] = true
	}

	// Get the list of this node's validator keys
	keys, err := w.GetAllPrivateKeys()
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error getting private validator keys: %w", err)
	}
	data.TotalCount = uint64(len(keys))

	// Find the ones that haven't been uploaded yet
	unregisteredKeys := []*eth2types.BLSPrivateKey{}
	newPubkeys := []beacon.ValidatorPubkey{}
	for _, key := range keys {
		pubkey := beacon.ValidatorPubkey(key.PublicKey().Marshal())
		_, exists := registeredPubkeyMap[pubkey]
		if !exists {
			unregisteredKeys = append(unregisteredKeys, key)
			newPubkeys = append(newPubkeys, pubkey)
		}
	}
	data.UnregisteredPubkeys = newPubkeys

	if len(unregisteredKeys) == 0 {
		return types.ResponseStatus_Success, nil
	}

	// Make sure validator has enough funds to pay for the deposit
	if !c.bypassBalanceCheck {
		balance, err := ec.BalanceAt(ctx, opts.From, nil)
		if err != nil {
			return types.ResponseStatus_Error, fmt.Errorf("error getting balance: %w", err)
		}
		data.Balance = balance

		totalCost := big.NewInt(int64(len(unregisteredKeys)))
		totalCost.Mul(totalCost, eth.EthToWei(0.01))
		data.RequiredBalance = totalCost

		data.SufficientBalance = (totalCost.Cmp(balance) < 0)
		if !data.SufficientBalance {
			return types.ResponseStatus_Success, nil
		}
	}

	// Get the deposit data for those pubkeys
	depositData, err := ddMgr.GenerateDepositData(unregisteredKeys)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error generating deposit data: %w", err)
	}

	// Serialize it
	bytes, err := json.Marshal(depositData)
	if err != nil {
		return types.ResponseStatus_Error, fmt.Errorf("error serializing deposit data: %w", err)
	}

	// Submit the upload
	response, err := nc.UploadDepositData(bytes)
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	data.ServerResponse = response
	return types.ResponseStatus_Success, nil
}
