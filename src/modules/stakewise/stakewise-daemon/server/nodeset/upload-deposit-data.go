package swnodeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"

	"github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/eth-utils/beacon"
	"github.com/nodeset-org/hyperdrive/daemon-utils/server"
	swapi "github.com/nodeset-org/hyperdrive/modules/stakewise/shared/api"
	"github.com/nodeset-org/hyperdrive/shared/utils/input"
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
	server.RegisterQuerylessGet[*nodesetUploadDepositDataContext, swapi.NodesetUploadDepositDataData](
		router, "upload-deposit-data", f, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type nodesetUploadDepositDataContext struct {
	handler            *NodesetHandler
	bypassBalanceCheck bool
}

func weiToEth(wei *big.Int) *big.Float {
	ether := big.NewFloat(1e18)

	weiFloat := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(weiFloat, ether)

	return ethValue
}

func (c *nodesetUploadDepositDataContext) PrepareData(data *swapi.NodesetUploadDepositDataData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	ddMgr := sp.GetDepositDataManager()
	nc := sp.GetNodesetClient()
	w := sp.GetWallet()
	ec := sp.GetEthClient()

	balance, err := ec.BalanceAt(context.Background(), opts.From, nil)
	if err != nil {
		return fmt.Errorf("error getting balance: %w", err)
	}

	// Get the list of registered validators
	registeredPubkeyMap := map[beacon.ValidatorPubkey]bool{}
	registeredPubkeys, err := nc.GetRegisteredValidators()
	if err != nil {
		return fmt.Errorf("error getting registered validators: %w", err)
	}
	for _, pubkey := range registeredPubkeys {
		registeredPubkeyMap[pubkey] = true
	}

	// Get the list of this node's validator keys
	keys, err := w.GetAllPrivateKeys()
	if err != nil {
		return fmt.Errorf("error getting private validator keys: %w", err)
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
	data.NewPubkeys = newPubkeys

	if len(unregisteredKeys) == 0 {
		return nil
	}

	// Make sure validator has enough funds to pay for the deposit
	if !c.bypassBalanceCheck {
		totalCost := new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(int64(len(unregisteredKeys))))
		if totalCost.Cmp(balance) > 0 {
			return fmt.Errorf("balance_check_failed: You're attempting to upload %v keys, but you only have %v ETH in your account. We recommend you have at least %v ETH", len(unregisteredKeys), weiToEth(balance), weiToEth(totalCost))
		}
	}

	// Get the deposit data for those pubkeys
	depositData, err := ddMgr.GenerateDepositData(unregisteredKeys)
	if err != nil {
		return fmt.Errorf("error generating deposit data: %w", err)
	}

	// Serialize it
	bytes, err := json.Marshal(depositData)
	if err != nil {
		return fmt.Errorf("error serializing deposit data: %w", err)
	}

	// Submit the upload
	response, err := nc.UploadDepositData(bytes)
	if err != nil {
		return err
	}
	data.ServerResponse = response
	return nil
}
