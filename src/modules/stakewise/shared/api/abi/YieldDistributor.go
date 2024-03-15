// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Claim is an auto generated low-level Go binding around an user-defined struct.
type Claim struct {
	Amount       *big.Int
	NumOperators *big.Int
}

// Reward is an auto generated low-level Go binding around an user-defined struct.
type Reward struct {
	Recipient common.Address
	Eth       *big.Int
}

// YieldDistributorMetaData contains all meta data concerning the YieldDistributor contract.
var YieldDistributorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"eth\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structReward\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"RewardDistributed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interval\",\"type\":\"uint256\"}],\"name\":\"WarningAlreadyClaimed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"}],\"name\":\"adminSweep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"claims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numOperators\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentIntervalGenesisTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dustAccrued\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInterval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getClaims\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numOperators\",\"type\":\"uint256\"}],\"internalType\":\"structClaim[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDirectory\",\"outputs\":[{\"internalType\":\"contractDirectory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rewardee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_startInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_endInterval\",\"type\":\"uint256\"}],\"name\":\"harvest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hasClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_directory\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxIntervalLengthSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_maxIntervalLengthSeconds\",\"type\":\"uint256\"}],\"name\":\"setMaxIntervalTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_k\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxValidators\",\"type\":\"uint256\"}],\"name\":\"setRewardIncentiveModel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalYieldAccrued\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"weth\",\"type\":\"uint256\"}],\"name\":\"wethReceived\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"yieldAccruedInInterval\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// YieldDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use YieldDistributorMetaData.ABI instead.
var YieldDistributorABI = YieldDistributorMetaData.ABI

// YieldDistributor is an auto generated Go binding around an Ethereum contract.
type YieldDistributor struct {
	YieldDistributorCaller     // Read-only binding to the contract
	YieldDistributorTransactor // Write-only binding to the contract
	YieldDistributorFilterer   // Log filterer for contract events
}

// YieldDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type YieldDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// YieldDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type YieldDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// YieldDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type YieldDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// YieldDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type YieldDistributorSession struct {
	Contract     *YieldDistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// YieldDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type YieldDistributorCallerSession struct {
	Contract *YieldDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// YieldDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type YieldDistributorTransactorSession struct {
	Contract     *YieldDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// YieldDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type YieldDistributorRaw struct {
	Contract *YieldDistributor // Generic contract binding to access the raw methods on
}

// YieldDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type YieldDistributorCallerRaw struct {
	Contract *YieldDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// YieldDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type YieldDistributorTransactorRaw struct {
	Contract *YieldDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewYieldDistributor creates a new instance of YieldDistributor, bound to a specific deployed contract.
func NewYieldDistributor(address common.Address, backend bind.ContractBackend) (*YieldDistributor, error) {
	contract, err := bindYieldDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &YieldDistributor{YieldDistributorCaller: YieldDistributorCaller{contract: contract}, YieldDistributorTransactor: YieldDistributorTransactor{contract: contract}, YieldDistributorFilterer: YieldDistributorFilterer{contract: contract}}, nil
}

// NewYieldDistributorCaller creates a new read-only instance of YieldDistributor, bound to a specific deployed contract.
func NewYieldDistributorCaller(address common.Address, caller bind.ContractCaller) (*YieldDistributorCaller, error) {
	contract, err := bindYieldDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &YieldDistributorCaller{contract: contract}, nil
}

// NewYieldDistributorTransactor creates a new write-only instance of YieldDistributor, bound to a specific deployed contract.
func NewYieldDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*YieldDistributorTransactor, error) {
	contract, err := bindYieldDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &YieldDistributorTransactor{contract: contract}, nil
}

// NewYieldDistributorFilterer creates a new log filterer instance of YieldDistributor, bound to a specific deployed contract.
func NewYieldDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*YieldDistributorFilterer, error) {
	contract, err := bindYieldDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &YieldDistributorFilterer{contract: contract}, nil
}

// bindYieldDistributor binds a generic wrapper to an already deployed contract.
func bindYieldDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := YieldDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_YieldDistributor *YieldDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _YieldDistributor.Contract.YieldDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_YieldDistributor *YieldDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldDistributor.Contract.YieldDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_YieldDistributor *YieldDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _YieldDistributor.Contract.YieldDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_YieldDistributor *YieldDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _YieldDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_YieldDistributor *YieldDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_YieldDistributor *YieldDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _YieldDistributor.Contract.contract.Transact(opts, method, params...)
}

// Claims is a free data retrieval call binding the contract method 0xa888c2cd.
//
// Solidity: function claims(uint256 ) view returns(uint256 amount, uint256 numOperators)
func (_YieldDistributor *YieldDistributorCaller) Claims(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount       *big.Int
	NumOperators *big.Int
}, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "claims", arg0)

	outstruct := new(struct {
		Amount       *big.Int
		NumOperators *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NumOperators = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Claims is a free data retrieval call binding the contract method 0xa888c2cd.
//
// Solidity: function claims(uint256 ) view returns(uint256 amount, uint256 numOperators)
func (_YieldDistributor *YieldDistributorSession) Claims(arg0 *big.Int) (struct {
	Amount       *big.Int
	NumOperators *big.Int
}, error) {
	return _YieldDistributor.Contract.Claims(&_YieldDistributor.CallOpts, arg0)
}

// Claims is a free data retrieval call binding the contract method 0xa888c2cd.
//
// Solidity: function claims(uint256 ) view returns(uint256 amount, uint256 numOperators)
func (_YieldDistributor *YieldDistributorCallerSession) Claims(arg0 *big.Int) (struct {
	Amount       *big.Int
	NumOperators *big.Int
}, error) {
	return _YieldDistributor.Contract.Claims(&_YieldDistributor.CallOpts, arg0)
}

// CurrentInterval is a free data retrieval call binding the contract method 0x363487bc.
//
// Solidity: function currentInterval() view returns(uint256)
func (_YieldDistributor *YieldDistributorCaller) CurrentInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "currentInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentInterval is a free data retrieval call binding the contract method 0x363487bc.
//
// Solidity: function currentInterval() view returns(uint256)
func (_YieldDistributor *YieldDistributorSession) CurrentInterval() (*big.Int, error) {
	return _YieldDistributor.Contract.CurrentInterval(&_YieldDistributor.CallOpts)
}

// CurrentInterval is a free data retrieval call binding the contract method 0x363487bc.
//
// Solidity: function currentInterval() view returns(uint256)
func (_YieldDistributor *YieldDistributorCallerSession) CurrentInterval() (*big.Int, error) {
	return _YieldDistributor.Contract.CurrentInterval(&_YieldDistributor.CallOpts)
}

// CurrentIntervalGenesisTime is a free data retrieval call binding the contract method 0x015c1af5.
//
// Solidity: function currentIntervalGenesisTime() view returns(uint256)
func (_YieldDistributor *YieldDistributorCaller) CurrentIntervalGenesisTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "currentIntervalGenesisTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentIntervalGenesisTime is a free data retrieval call binding the contract method 0x015c1af5.
//
// Solidity: function currentIntervalGenesisTime() view returns(uint256)
func (_YieldDistributor *YieldDistributorSession) CurrentIntervalGenesisTime() (*big.Int, error) {
	return _YieldDistributor.Contract.CurrentIntervalGenesisTime(&_YieldDistributor.CallOpts)
}

// CurrentIntervalGenesisTime is a free data retrieval call binding the contract method 0x015c1af5.
//
// Solidity: function currentIntervalGenesisTime() view returns(uint256)
func (_YieldDistributor *YieldDistributorCallerSession) CurrentIntervalGenesisTime() (*big.Int, error) {
	return _YieldDistributor.Contract.CurrentIntervalGenesisTime(&_YieldDistributor.CallOpts)
}

// DustAccrued is a free data retrieval call binding the contract method 0x9adb3844.
//
// Solidity: function dustAccrued() view returns(uint256)
func (_YieldDistributor *YieldDistributorCaller) DustAccrued(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "dustAccrued")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DustAccrued is a free data retrieval call binding the contract method 0x9adb3844.
//
// Solidity: function dustAccrued() view returns(uint256)
func (_YieldDistributor *YieldDistributorSession) DustAccrued() (*big.Int, error) {
	return _YieldDistributor.Contract.DustAccrued(&_YieldDistributor.CallOpts)
}

// DustAccrued is a free data retrieval call binding the contract method 0x9adb3844.
//
// Solidity: function dustAccrued() view returns(uint256)
func (_YieldDistributor *YieldDistributorCallerSession) DustAccrued() (*big.Int, error) {
	return _YieldDistributor.Contract.DustAccrued(&_YieldDistributor.CallOpts)
}

// GetClaims is a free data retrieval call binding the contract method 0xc52822f8.
//
// Solidity: function getClaims() view returns((uint256,uint256)[])
func (_YieldDistributor *YieldDistributorCaller) GetClaims(opts *bind.CallOpts) ([]Claim, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "getClaims")

	if err != nil {
		return *new([]Claim), err
	}

	out0 := *abi.ConvertType(out[0], new([]Claim)).(*[]Claim)

	return out0, err

}

// GetClaims is a free data retrieval call binding the contract method 0xc52822f8.
//
// Solidity: function getClaims() view returns((uint256,uint256)[])
func (_YieldDistributor *YieldDistributorSession) GetClaims() ([]Claim, error) {
	return _YieldDistributor.Contract.GetClaims(&_YieldDistributor.CallOpts)
}

// GetClaims is a free data retrieval call binding the contract method 0xc52822f8.
//
// Solidity: function getClaims() view returns((uint256,uint256)[])
func (_YieldDistributor *YieldDistributorCallerSession) GetClaims() ([]Claim, error) {
	return _YieldDistributor.Contract.GetClaims(&_YieldDistributor.CallOpts)
}

// GetDirectory is a free data retrieval call binding the contract method 0x76247776.
//
// Solidity: function getDirectory() view returns(address)
func (_YieldDistributor *YieldDistributorCaller) GetDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "getDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDirectory is a free data retrieval call binding the contract method 0x76247776.
//
// Solidity: function getDirectory() view returns(address)
func (_YieldDistributor *YieldDistributorSession) GetDirectory() (common.Address, error) {
	return _YieldDistributor.Contract.GetDirectory(&_YieldDistributor.CallOpts)
}

// GetDirectory is a free data retrieval call binding the contract method 0x76247776.
//
// Solidity: function getDirectory() view returns(address)
func (_YieldDistributor *YieldDistributorCallerSession) GetDirectory() (common.Address, error) {
	return _YieldDistributor.Contract.GetDirectory(&_YieldDistributor.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_YieldDistributor *YieldDistributorCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_YieldDistributor *YieldDistributorSession) GetImplementation() (common.Address, error) {
	return _YieldDistributor.Contract.GetImplementation(&_YieldDistributor.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_YieldDistributor *YieldDistributorCallerSession) GetImplementation() (common.Address, error) {
	return _YieldDistributor.Contract.GetImplementation(&_YieldDistributor.CallOpts)
}

// HasClaimed is a free data retrieval call binding the contract method 0xb2931096.
//
// Solidity: function hasClaimed(address , uint256 ) view returns(bool)
func (_YieldDistributor *YieldDistributorCaller) HasClaimed(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "hasClaimed", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasClaimed is a free data retrieval call binding the contract method 0xb2931096.
//
// Solidity: function hasClaimed(address , uint256 ) view returns(bool)
func (_YieldDistributor *YieldDistributorSession) HasClaimed(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _YieldDistributor.Contract.HasClaimed(&_YieldDistributor.CallOpts, arg0, arg1)
}

// HasClaimed is a free data retrieval call binding the contract method 0xb2931096.
//
// Solidity: function hasClaimed(address , uint256 ) view returns(bool)
func (_YieldDistributor *YieldDistributorCallerSession) HasClaimed(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _YieldDistributor.Contract.HasClaimed(&_YieldDistributor.CallOpts, arg0, arg1)
}

// MaxIntervalLengthSeconds is a free data retrieval call binding the contract method 0x6b8680cb.
//
// Solidity: function maxIntervalLengthSeconds() view returns(uint256)
func (_YieldDistributor *YieldDistributorCaller) MaxIntervalLengthSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "maxIntervalLengthSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxIntervalLengthSeconds is a free data retrieval call binding the contract method 0x6b8680cb.
//
// Solidity: function maxIntervalLengthSeconds() view returns(uint256)
func (_YieldDistributor *YieldDistributorSession) MaxIntervalLengthSeconds() (*big.Int, error) {
	return _YieldDistributor.Contract.MaxIntervalLengthSeconds(&_YieldDistributor.CallOpts)
}

// MaxIntervalLengthSeconds is a free data retrieval call binding the contract method 0x6b8680cb.
//
// Solidity: function maxIntervalLengthSeconds() view returns(uint256)
func (_YieldDistributor *YieldDistributorCallerSession) MaxIntervalLengthSeconds() (*big.Int, error) {
	return _YieldDistributor.Contract.MaxIntervalLengthSeconds(&_YieldDistributor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_YieldDistributor *YieldDistributorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_YieldDistributor *YieldDistributorSession) ProxiableUUID() ([32]byte, error) {
	return _YieldDistributor.Contract.ProxiableUUID(&_YieldDistributor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_YieldDistributor *YieldDistributorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _YieldDistributor.Contract.ProxiableUUID(&_YieldDistributor.CallOpts)
}

// TotalYieldAccrued is a free data retrieval call binding the contract method 0xee6c3bf2.
//
// Solidity: function totalYieldAccrued() view returns(uint256)
func (_YieldDistributor *YieldDistributorCaller) TotalYieldAccrued(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "totalYieldAccrued")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalYieldAccrued is a free data retrieval call binding the contract method 0xee6c3bf2.
//
// Solidity: function totalYieldAccrued() view returns(uint256)
func (_YieldDistributor *YieldDistributorSession) TotalYieldAccrued() (*big.Int, error) {
	return _YieldDistributor.Contract.TotalYieldAccrued(&_YieldDistributor.CallOpts)
}

// TotalYieldAccrued is a free data retrieval call binding the contract method 0xee6c3bf2.
//
// Solidity: function totalYieldAccrued() view returns(uint256)
func (_YieldDistributor *YieldDistributorCallerSession) TotalYieldAccrued() (*big.Int, error) {
	return _YieldDistributor.Contract.TotalYieldAccrued(&_YieldDistributor.CallOpts)
}

// YieldAccruedInInterval is a free data retrieval call binding the contract method 0x3d1a4637.
//
// Solidity: function yieldAccruedInInterval() view returns(uint256)
func (_YieldDistributor *YieldDistributorCaller) YieldAccruedInInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _YieldDistributor.contract.Call(opts, &out, "yieldAccruedInInterval")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// YieldAccruedInInterval is a free data retrieval call binding the contract method 0x3d1a4637.
//
// Solidity: function yieldAccruedInInterval() view returns(uint256)
func (_YieldDistributor *YieldDistributorSession) YieldAccruedInInterval() (*big.Int, error) {
	return _YieldDistributor.Contract.YieldAccruedInInterval(&_YieldDistributor.CallOpts)
}

// YieldAccruedInInterval is a free data retrieval call binding the contract method 0x3d1a4637.
//
// Solidity: function yieldAccruedInInterval() view returns(uint256)
func (_YieldDistributor *YieldDistributorCallerSession) YieldAccruedInInterval() (*big.Int, error) {
	return _YieldDistributor.Contract.YieldAccruedInInterval(&_YieldDistributor.CallOpts)
}

// AdminSweep is a paid mutator transaction binding the contract method 0xb7bf46af.
//
// Solidity: function adminSweep(address treasury) returns()
func (_YieldDistributor *YieldDistributorTransactor) AdminSweep(opts *bind.TransactOpts, treasury common.Address) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "adminSweep", treasury)
}

// AdminSweep is a paid mutator transaction binding the contract method 0xb7bf46af.
//
// Solidity: function adminSweep(address treasury) returns()
func (_YieldDistributor *YieldDistributorSession) AdminSweep(treasury common.Address) (*types.Transaction, error) {
	return _YieldDistributor.Contract.AdminSweep(&_YieldDistributor.TransactOpts, treasury)
}

// AdminSweep is a paid mutator transaction binding the contract method 0xb7bf46af.
//
// Solidity: function adminSweep(address treasury) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) AdminSweep(treasury common.Address) (*types.Transaction, error) {
	return _YieldDistributor.Contract.AdminSweep(&_YieldDistributor.TransactOpts, treasury)
}

// FinalizeInterval is a paid mutator transaction binding the contract method 0x27c448c6.
//
// Solidity: function finalizeInterval() returns()
func (_YieldDistributor *YieldDistributorTransactor) FinalizeInterval(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "finalizeInterval")
}

// FinalizeInterval is a paid mutator transaction binding the contract method 0x27c448c6.
//
// Solidity: function finalizeInterval() returns()
func (_YieldDistributor *YieldDistributorSession) FinalizeInterval() (*types.Transaction, error) {
	return _YieldDistributor.Contract.FinalizeInterval(&_YieldDistributor.TransactOpts)
}

// FinalizeInterval is a paid mutator transaction binding the contract method 0x27c448c6.
//
// Solidity: function finalizeInterval() returns()
func (_YieldDistributor *YieldDistributorTransactorSession) FinalizeInterval() (*types.Transaction, error) {
	return _YieldDistributor.Contract.FinalizeInterval(&_YieldDistributor.TransactOpts)
}

// Harvest is a paid mutator transaction binding the contract method 0xf7588701.
//
// Solidity: function harvest(address _rewardee, uint256 _startInterval, uint256 _endInterval) returns()
func (_YieldDistributor *YieldDistributorTransactor) Harvest(opts *bind.TransactOpts, _rewardee common.Address, _startInterval *big.Int, _endInterval *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "harvest", _rewardee, _startInterval, _endInterval)
}

// Harvest is a paid mutator transaction binding the contract method 0xf7588701.
//
// Solidity: function harvest(address _rewardee, uint256 _startInterval, uint256 _endInterval) returns()
func (_YieldDistributor *YieldDistributorSession) Harvest(_rewardee common.Address, _startInterval *big.Int, _endInterval *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.Harvest(&_YieldDistributor.TransactOpts, _rewardee, _startInterval, _endInterval)
}

// Harvest is a paid mutator transaction binding the contract method 0xf7588701.
//
// Solidity: function harvest(address _rewardee, uint256 _startInterval, uint256 _endInterval) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) Harvest(_rewardee common.Address, _startInterval *big.Int, _endInterval *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.Harvest(&_YieldDistributor.TransactOpts, _rewardee, _startInterval, _endInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _directory) returns()
func (_YieldDistributor *YieldDistributorTransactor) Initialize(opts *bind.TransactOpts, _directory common.Address) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "initialize", _directory)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _directory) returns()
func (_YieldDistributor *YieldDistributorSession) Initialize(_directory common.Address) (*types.Transaction, error) {
	return _YieldDistributor.Contract.Initialize(&_YieldDistributor.TransactOpts, _directory)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _directory) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) Initialize(_directory common.Address) (*types.Transaction, error) {
	return _YieldDistributor.Contract.Initialize(&_YieldDistributor.TransactOpts, _directory)
}

// SetMaxIntervalTime is a paid mutator transaction binding the contract method 0x02ef103b.
//
// Solidity: function setMaxIntervalTime(uint256 _maxIntervalLengthSeconds) returns()
func (_YieldDistributor *YieldDistributorTransactor) SetMaxIntervalTime(opts *bind.TransactOpts, _maxIntervalLengthSeconds *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "setMaxIntervalTime", _maxIntervalLengthSeconds)
}

// SetMaxIntervalTime is a paid mutator transaction binding the contract method 0x02ef103b.
//
// Solidity: function setMaxIntervalTime(uint256 _maxIntervalLengthSeconds) returns()
func (_YieldDistributor *YieldDistributorSession) SetMaxIntervalTime(_maxIntervalLengthSeconds *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.SetMaxIntervalTime(&_YieldDistributor.TransactOpts, _maxIntervalLengthSeconds)
}

// SetMaxIntervalTime is a paid mutator transaction binding the contract method 0x02ef103b.
//
// Solidity: function setMaxIntervalTime(uint256 _maxIntervalLengthSeconds) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) SetMaxIntervalTime(_maxIntervalLengthSeconds *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.SetMaxIntervalTime(&_YieldDistributor.TransactOpts, _maxIntervalLengthSeconds)
}

// SetRewardIncentiveModel is a paid mutator transaction binding the contract method 0xb88a4441.
//
// Solidity: function setRewardIncentiveModel(uint256 _k, uint256 _maxValidators) returns()
func (_YieldDistributor *YieldDistributorTransactor) SetRewardIncentiveModel(opts *bind.TransactOpts, _k *big.Int, _maxValidators *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "setRewardIncentiveModel", _k, _maxValidators)
}

// SetRewardIncentiveModel is a paid mutator transaction binding the contract method 0xb88a4441.
//
// Solidity: function setRewardIncentiveModel(uint256 _k, uint256 _maxValidators) returns()
func (_YieldDistributor *YieldDistributorSession) SetRewardIncentiveModel(_k *big.Int, _maxValidators *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.SetRewardIncentiveModel(&_YieldDistributor.TransactOpts, _k, _maxValidators)
}

// SetRewardIncentiveModel is a paid mutator transaction binding the contract method 0xb88a4441.
//
// Solidity: function setRewardIncentiveModel(uint256 _k, uint256 _maxValidators) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) SetRewardIncentiveModel(_k *big.Int, _maxValidators *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.SetRewardIncentiveModel(&_YieldDistributor.TransactOpts, _k, _maxValidators)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_YieldDistributor *YieldDistributorTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_YieldDistributor *YieldDistributorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _YieldDistributor.Contract.UpgradeTo(&_YieldDistributor.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _YieldDistributor.Contract.UpgradeTo(&_YieldDistributor.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_YieldDistributor *YieldDistributorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_YieldDistributor *YieldDistributorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _YieldDistributor.Contract.UpgradeToAndCall(&_YieldDistributor.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_YieldDistributor *YieldDistributorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _YieldDistributor.Contract.UpgradeToAndCall(&_YieldDistributor.TransactOpts, newImplementation, data)
}

// WethReceived is a paid mutator transaction binding the contract method 0x489380f8.
//
// Solidity: function wethReceived(uint256 weth) returns()
func (_YieldDistributor *YieldDistributorTransactor) WethReceived(opts *bind.TransactOpts, weth *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.contract.Transact(opts, "wethReceived", weth)
}

// WethReceived is a paid mutator transaction binding the contract method 0x489380f8.
//
// Solidity: function wethReceived(uint256 weth) returns()
func (_YieldDistributor *YieldDistributorSession) WethReceived(weth *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.WethReceived(&_YieldDistributor.TransactOpts, weth)
}

// WethReceived is a paid mutator transaction binding the contract method 0x489380f8.
//
// Solidity: function wethReceived(uint256 weth) returns()
func (_YieldDistributor *YieldDistributorTransactorSession) WethReceived(weth *big.Int) (*types.Transaction, error) {
	return _YieldDistributor.Contract.WethReceived(&_YieldDistributor.TransactOpts, weth)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_YieldDistributor *YieldDistributorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _YieldDistributor.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_YieldDistributor *YieldDistributorSession) Receive() (*types.Transaction, error) {
	return _YieldDistributor.Contract.Receive(&_YieldDistributor.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_YieldDistributor *YieldDistributorTransactorSession) Receive() (*types.Transaction, error) {
	return _YieldDistributor.Contract.Receive(&_YieldDistributor.TransactOpts)
}

// YieldDistributorAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the YieldDistributor contract.
type YieldDistributorAdminChangedIterator struct {
	Event *YieldDistributorAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldDistributorAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldDistributorAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldDistributorAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldDistributorAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldDistributorAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldDistributorAdminChanged represents a AdminChanged event raised by the YieldDistributor contract.
type YieldDistributorAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_YieldDistributor *YieldDistributorFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*YieldDistributorAdminChangedIterator, error) {

	logs, sub, err := _YieldDistributor.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &YieldDistributorAdminChangedIterator{contract: _YieldDistributor.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_YieldDistributor *YieldDistributorFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *YieldDistributorAdminChanged) (event.Subscription, error) {

	logs, sub, err := _YieldDistributor.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldDistributorAdminChanged)
				if err := _YieldDistributor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_YieldDistributor *YieldDistributorFilterer) ParseAdminChanged(log types.Log) (*YieldDistributorAdminChanged, error) {
	event := new(YieldDistributorAdminChanged)
	if err := _YieldDistributor.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldDistributorBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the YieldDistributor contract.
type YieldDistributorBeaconUpgradedIterator struct {
	Event *YieldDistributorBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldDistributorBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldDistributorBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldDistributorBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldDistributorBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldDistributorBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldDistributorBeaconUpgraded represents a BeaconUpgraded event raised by the YieldDistributor contract.
type YieldDistributorBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_YieldDistributor *YieldDistributorFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*YieldDistributorBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _YieldDistributor.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &YieldDistributorBeaconUpgradedIterator{contract: _YieldDistributor.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_YieldDistributor *YieldDistributorFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *YieldDistributorBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _YieldDistributor.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldDistributorBeaconUpgraded)
				if err := _YieldDistributor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_YieldDistributor *YieldDistributorFilterer) ParseBeaconUpgraded(log types.Log) (*YieldDistributorBeaconUpgraded, error) {
	event := new(YieldDistributorBeaconUpgraded)
	if err := _YieldDistributor.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldDistributorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the YieldDistributor contract.
type YieldDistributorInitializedIterator struct {
	Event *YieldDistributorInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldDistributorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldDistributorInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldDistributorInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldDistributorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldDistributorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldDistributorInitialized represents a Initialized event raised by the YieldDistributor contract.
type YieldDistributorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_YieldDistributor *YieldDistributorFilterer) FilterInitialized(opts *bind.FilterOpts) (*YieldDistributorInitializedIterator, error) {

	logs, sub, err := _YieldDistributor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &YieldDistributorInitializedIterator{contract: _YieldDistributor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_YieldDistributor *YieldDistributorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *YieldDistributorInitialized) (event.Subscription, error) {

	logs, sub, err := _YieldDistributor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldDistributorInitialized)
				if err := _YieldDistributor.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_YieldDistributor *YieldDistributorFilterer) ParseInitialized(log types.Log) (*YieldDistributorInitialized, error) {
	event := new(YieldDistributorInitialized)
	if err := _YieldDistributor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldDistributorRewardDistributedIterator is returned from FilterRewardDistributed and is used to iterate over the raw logs and unpacked data for RewardDistributed events raised by the YieldDistributor contract.
type YieldDistributorRewardDistributedIterator struct {
	Event *YieldDistributorRewardDistributed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldDistributorRewardDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldDistributorRewardDistributed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldDistributorRewardDistributed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldDistributorRewardDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldDistributorRewardDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldDistributorRewardDistributed represents a RewardDistributed event raised by the YieldDistributor contract.
type YieldDistributorRewardDistributed struct {
	Arg0 Reward
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRewardDistributed is a free log retrieval operation binding the contract event 0x035694c65c4febca8d92b53152411f8c2e10489ba0bdaa40844885002b5a52cb.
//
// Solidity: event RewardDistributed((address,uint256) arg0)
func (_YieldDistributor *YieldDistributorFilterer) FilterRewardDistributed(opts *bind.FilterOpts) (*YieldDistributorRewardDistributedIterator, error) {

	logs, sub, err := _YieldDistributor.contract.FilterLogs(opts, "RewardDistributed")
	if err != nil {
		return nil, err
	}
	return &YieldDistributorRewardDistributedIterator{contract: _YieldDistributor.contract, event: "RewardDistributed", logs: logs, sub: sub}, nil
}

// WatchRewardDistributed is a free log subscription operation binding the contract event 0x035694c65c4febca8d92b53152411f8c2e10489ba0bdaa40844885002b5a52cb.
//
// Solidity: event RewardDistributed((address,uint256) arg0)
func (_YieldDistributor *YieldDistributorFilterer) WatchRewardDistributed(opts *bind.WatchOpts, sink chan<- *YieldDistributorRewardDistributed) (event.Subscription, error) {

	logs, sub, err := _YieldDistributor.contract.WatchLogs(opts, "RewardDistributed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldDistributorRewardDistributed)
				if err := _YieldDistributor.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardDistributed is a log parse operation binding the contract event 0x035694c65c4febca8d92b53152411f8c2e10489ba0bdaa40844885002b5a52cb.
//
// Solidity: event RewardDistributed((address,uint256) arg0)
func (_YieldDistributor *YieldDistributorFilterer) ParseRewardDistributed(log types.Log) (*YieldDistributorRewardDistributed, error) {
	event := new(YieldDistributorRewardDistributed)
	if err := _YieldDistributor.contract.UnpackLog(event, "RewardDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldDistributorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the YieldDistributor contract.
type YieldDistributorUpgradedIterator struct {
	Event *YieldDistributorUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldDistributorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldDistributorUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldDistributorUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldDistributorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldDistributorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldDistributorUpgraded represents a Upgraded event raised by the YieldDistributor contract.
type YieldDistributorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_YieldDistributor *YieldDistributorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*YieldDistributorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _YieldDistributor.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &YieldDistributorUpgradedIterator{contract: _YieldDistributor.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_YieldDistributor *YieldDistributorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *YieldDistributorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _YieldDistributor.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldDistributorUpgraded)
				if err := _YieldDistributor.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_YieldDistributor *YieldDistributorFilterer) ParseUpgraded(log types.Log) (*YieldDistributorUpgraded, error) {
	event := new(YieldDistributorUpgraded)
	if err := _YieldDistributor.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// YieldDistributorWarningAlreadyClaimedIterator is returned from FilterWarningAlreadyClaimed and is used to iterate over the raw logs and unpacked data for WarningAlreadyClaimed events raised by the YieldDistributor contract.
type YieldDistributorWarningAlreadyClaimedIterator struct {
	Event *YieldDistributorWarningAlreadyClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *YieldDistributorWarningAlreadyClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(YieldDistributorWarningAlreadyClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(YieldDistributorWarningAlreadyClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *YieldDistributorWarningAlreadyClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *YieldDistributorWarningAlreadyClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// YieldDistributorWarningAlreadyClaimed represents a WarningAlreadyClaimed event raised by the YieldDistributor contract.
type YieldDistributorWarningAlreadyClaimed struct {
	Operator common.Address
	Interval *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWarningAlreadyClaimed is a free log retrieval operation binding the contract event 0x568de2b5ebb1774db867a49e2c54e71438d11cb7d48ec787f98b15dd117dc34f.
//
// Solidity: event WarningAlreadyClaimed(address operator, uint256 interval)
func (_YieldDistributor *YieldDistributorFilterer) FilterWarningAlreadyClaimed(opts *bind.FilterOpts) (*YieldDistributorWarningAlreadyClaimedIterator, error) {

	logs, sub, err := _YieldDistributor.contract.FilterLogs(opts, "WarningAlreadyClaimed")
	if err != nil {
		return nil, err
	}
	return &YieldDistributorWarningAlreadyClaimedIterator{contract: _YieldDistributor.contract, event: "WarningAlreadyClaimed", logs: logs, sub: sub}, nil
}

// WatchWarningAlreadyClaimed is a free log subscription operation binding the contract event 0x568de2b5ebb1774db867a49e2c54e71438d11cb7d48ec787f98b15dd117dc34f.
//
// Solidity: event WarningAlreadyClaimed(address operator, uint256 interval)
func (_YieldDistributor *YieldDistributorFilterer) WatchWarningAlreadyClaimed(opts *bind.WatchOpts, sink chan<- *YieldDistributorWarningAlreadyClaimed) (event.Subscription, error) {

	logs, sub, err := _YieldDistributor.contract.WatchLogs(opts, "WarningAlreadyClaimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(YieldDistributorWarningAlreadyClaimed)
				if err := _YieldDistributor.contract.UnpackLog(event, "WarningAlreadyClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWarningAlreadyClaimed is a log parse operation binding the contract event 0x568de2b5ebb1774db867a49e2c54e71438d11cb7d48ec787f98b15dd117dc34f.
//
// Solidity: event WarningAlreadyClaimed(address operator, uint256 interval)
func (_YieldDistributor *YieldDistributorFilterer) ParseWarningAlreadyClaimed(log types.Log) (*YieldDistributorWarningAlreadyClaimed, error) {
	event := new(YieldDistributorWarningAlreadyClaimed)
	if err := _YieldDistributor.contract.UnpackLog(event, "WarningAlreadyClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
