package with_ns_registered

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
)

const (
	whitelistAbiString string = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"operationStartTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"intervalStart\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"operatorController\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"struct Operator\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldController\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newController\",\"type\":\"address\"}],\"name\":\"OperatorControllerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"OperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"OperatorsAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"OperatorsRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_sig\",\"type\":\"bytes\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"_sig\",\"type\":\"bytes[]\"}],\"name\":\"addOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDirectory\",\"outputs\":[{\"internalType\":\"contract Directory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"getIsAddressInWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"getNumberOfValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getOperatorAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"getOperatorAtAddress\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"operationStartTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"intervalStart\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"operatorController\",\"type\":\"address\"}],\"internalType\":\"struct Operator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"directoryAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"directoryAddress\",\"type\":\"address\"}],\"name\":\"initializeWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"nodeIndexMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nodeMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"operationStartTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"intervalStart\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"operatorController\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numOperators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorControllerToNodeMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeOperator\",\"type\":\"address\"}],\"name\":\"registerNewValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeOperator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"removeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"reverseNodeIndexMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"setOperatorController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sigsUsed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"
)

// ABI cache
var whitelistAbi abi.ABI
var whitelistOnce sync.Once

type Whitelist struct {
	contract *eth.Contract
	txMgr    *eth.TransactionManager
}

// Create a new Whitelist instance
func NewWhitelist(address common.Address, ec eth.IExecutionClient, txMgr *eth.TransactionManager) (*Whitelist, error) {
	// Parse the ABI
	var err error
	whitelistOnce.Do(func() {
		var parsedAbi abi.ABI
		parsedAbi, err = abi.JSON(strings.NewReader(whitelistAbiString))
		if err == nil {
			whitelistAbi = parsedAbi
		}
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing SplitWarehouse ABI: %w", err)
	}

	// Create the contract
	contract := &eth.Contract{
		ContractImpl: bind.NewBoundContract(address, whitelistAbi, ec, ec, ec),
		Address:      address,
		ABI:          &whitelistAbi,
	}

	return &Whitelist{
		contract: contract,
		txMgr:    txMgr,
	}, nil
}

// =============
// === Calls ===
// =============

func (c *Whitelist) IsAddressInWhitelist(mc *batch.MultiCaller, out *bool, account common.Address) {
	eth.AddCallToMulticaller(mc, c.contract, out, "getIsAddressInWhitelist", account)
}

// ====================
// === Transactions ===
// ====================

func (c *Whitelist) AddOperator(address common.Address, signature []byte, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.contract, "addOperator", opts, address, signature)
}
