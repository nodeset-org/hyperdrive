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

// SplitMainMetaData contains all meta data concerning the SplitMain contract.
var SplitMainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Create2Error\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CreateError\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newController\",\"type\":\"address\"}],\"name\":\"InvalidNewController\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"accountsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"allocationsLength\",\"type\":\"uint256\"}],\"name\":\"InvalidSplit__AccountsAndAllocationsMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"InvalidSplit__AccountsOutOfOrder\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"InvalidSplit__AllocationMustBePositive\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"allocationsSum\",\"type\":\"uint32\"}],\"name\":\"InvalidSplit__InvalidAllocationsSum\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"}],\"name\":\"InvalidSplit__InvalidDistributorFee\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"InvalidSplit__InvalidHash\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"accountsLength\",\"type\":\"uint256\"}],\"name\":\"InvalidSplit__TooFewAccounts\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"CancelControlTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousController\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newController\",\"type\":\"address\"}],\"name\":\"ControlTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"CreateSplit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"distributorAddress\",\"type\":\"address\"}],\"name\":\"DistributeERC20\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"distributorAddress\",\"type\":\"address\"}],\"name\":\"DistributeETH\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newPotentialController\",\"type\":\"address\"}],\"name\":\"InitiateControlTransfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"UpdateSplit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"contractERC20[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenAmounts\",\"type\":\"uint256[]\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"PERCENTAGE_SCALE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"acceptControl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"cancelControlTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"createSplit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"internalType\":\"contractERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"distributorAddress\",\"type\":\"address\"}],\"name\":\"distributeERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"distributorAddress\",\"type\":\"address\"}],\"name\":\"distributeETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"getController\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"contractERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getERC20Balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getETHBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"getHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"getNewPotentialController\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"name\":\"makeSplitImmutable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"}],\"name\":\"predictImmutableSplitAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newController\",\"type\":\"address\"}],\"name\":\"transferControl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"internalType\":\"contractERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"distributorAddress\",\"type\":\"address\"}],\"name\":\"updateAndDistributeERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"distributorAddress\",\"type\":\"address\"}],\"name\":\"updateAndDistributeETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint32[]\",\"name\":\"percentAllocations\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32\",\"name\":\"distributorFee\",\"type\":\"uint32\"}],\"name\":\"updateSplit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"walletImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"withdrawETH\",\"type\":\"uint256\"},{\"internalType\":\"contractERC20[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// SplitMainABI is the input ABI used to generate the binding from.
// Deprecated: Use SplitMainMetaData.ABI instead.
var SplitMainABI = SplitMainMetaData.ABI

// SplitMain is an auto generated Go binding around an Ethereum contract.
type SplitMain struct {
	SplitMainCaller     // Read-only binding to the contract
	SplitMainTransactor // Write-only binding to the contract
	SplitMainFilterer   // Log filterer for contract events
}

// SplitMainCaller is an auto generated read-only Go binding around an Ethereum contract.
type SplitMainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SplitMainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SplitMainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SplitMainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SplitMainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SplitMainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SplitMainSession struct {
	Contract     *SplitMain        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SplitMainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SplitMainCallerSession struct {
	Contract *SplitMainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SplitMainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SplitMainTransactorSession struct {
	Contract     *SplitMainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SplitMainRaw is an auto generated low-level Go binding around an Ethereum contract.
type SplitMainRaw struct {
	Contract *SplitMain // Generic contract binding to access the raw methods on
}

// SplitMainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SplitMainCallerRaw struct {
	Contract *SplitMainCaller // Generic read-only contract binding to access the raw methods on
}

// SplitMainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SplitMainTransactorRaw struct {
	Contract *SplitMainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSplitMain creates a new instance of SplitMain, bound to a specific deployed contract.
func NewSplitMain(address common.Address, backend bind.ContractBackend) (*SplitMain, error) {
	contract, err := bindSplitMain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SplitMain{SplitMainCaller: SplitMainCaller{contract: contract}, SplitMainTransactor: SplitMainTransactor{contract: contract}, SplitMainFilterer: SplitMainFilterer{contract: contract}}, nil
}

// NewSplitMainCaller creates a new read-only instance of SplitMain, bound to a specific deployed contract.
func NewSplitMainCaller(address common.Address, caller bind.ContractCaller) (*SplitMainCaller, error) {
	contract, err := bindSplitMain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SplitMainCaller{contract: contract}, nil
}

// NewSplitMainTransactor creates a new write-only instance of SplitMain, bound to a specific deployed contract.
func NewSplitMainTransactor(address common.Address, transactor bind.ContractTransactor) (*SplitMainTransactor, error) {
	contract, err := bindSplitMain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SplitMainTransactor{contract: contract}, nil
}

// NewSplitMainFilterer creates a new log filterer instance of SplitMain, bound to a specific deployed contract.
func NewSplitMainFilterer(address common.Address, filterer bind.ContractFilterer) (*SplitMainFilterer, error) {
	contract, err := bindSplitMain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SplitMainFilterer{contract: contract}, nil
}

// bindSplitMain binds a generic wrapper to an already deployed contract.
func bindSplitMain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SplitMainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SplitMain *SplitMainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SplitMain.Contract.SplitMainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SplitMain *SplitMainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SplitMain.Contract.SplitMainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SplitMain *SplitMainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SplitMain.Contract.SplitMainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SplitMain *SplitMainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SplitMain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SplitMain *SplitMainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SplitMain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SplitMain *SplitMainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SplitMain.Contract.contract.Transact(opts, method, params...)
}

// PERCENTAGESCALE is a free data retrieval call binding the contract method 0x3f26479e.
//
// Solidity: function PERCENTAGE_SCALE() view returns(uint256)
func (_SplitMain *SplitMainCaller) PERCENTAGESCALE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "PERCENTAGE_SCALE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PERCENTAGESCALE is a free data retrieval call binding the contract method 0x3f26479e.
//
// Solidity: function PERCENTAGE_SCALE() view returns(uint256)
func (_SplitMain *SplitMainSession) PERCENTAGESCALE() (*big.Int, error) {
	return _SplitMain.Contract.PERCENTAGESCALE(&_SplitMain.CallOpts)
}

// PERCENTAGESCALE is a free data retrieval call binding the contract method 0x3f26479e.
//
// Solidity: function PERCENTAGE_SCALE() view returns(uint256)
func (_SplitMain *SplitMainCallerSession) PERCENTAGESCALE() (*big.Int, error) {
	return _SplitMain.Contract.PERCENTAGESCALE(&_SplitMain.CallOpts)
}

// GetController is a free data retrieval call binding the contract method 0x88c662aa.
//
// Solidity: function getController(address split) view returns(address)
func (_SplitMain *SplitMainCaller) GetController(opts *bind.CallOpts, split common.Address) (common.Address, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "getController", split)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetController is a free data retrieval call binding the contract method 0x88c662aa.
//
// Solidity: function getController(address split) view returns(address)
func (_SplitMain *SplitMainSession) GetController(split common.Address) (common.Address, error) {
	return _SplitMain.Contract.GetController(&_SplitMain.CallOpts, split)
}

// GetController is a free data retrieval call binding the contract method 0x88c662aa.
//
// Solidity: function getController(address split) view returns(address)
func (_SplitMain *SplitMainCallerSession) GetController(split common.Address) (common.Address, error) {
	return _SplitMain.Contract.GetController(&_SplitMain.CallOpts, split)
}

// GetERC20Balance is a free data retrieval call binding the contract method 0xc3a8962c.
//
// Solidity: function getERC20Balance(address account, address token) view returns(uint256)
func (_SplitMain *SplitMainCaller) GetERC20Balance(opts *bind.CallOpts, account common.Address, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "getERC20Balance", account, token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetERC20Balance is a free data retrieval call binding the contract method 0xc3a8962c.
//
// Solidity: function getERC20Balance(address account, address token) view returns(uint256)
func (_SplitMain *SplitMainSession) GetERC20Balance(account common.Address, token common.Address) (*big.Int, error) {
	return _SplitMain.Contract.GetERC20Balance(&_SplitMain.CallOpts, account, token)
}

// GetERC20Balance is a free data retrieval call binding the contract method 0xc3a8962c.
//
// Solidity: function getERC20Balance(address account, address token) view returns(uint256)
func (_SplitMain *SplitMainCallerSession) GetERC20Balance(account common.Address, token common.Address) (*big.Int, error) {
	return _SplitMain.Contract.GetERC20Balance(&_SplitMain.CallOpts, account, token)
}

// GetETHBalance is a free data retrieval call binding the contract method 0x3bb66a7b.
//
// Solidity: function getETHBalance(address account) view returns(uint256)
func (_SplitMain *SplitMainCaller) GetETHBalance(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "getETHBalance", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetETHBalance is a free data retrieval call binding the contract method 0x3bb66a7b.
//
// Solidity: function getETHBalance(address account) view returns(uint256)
func (_SplitMain *SplitMainSession) GetETHBalance(account common.Address) (*big.Int, error) {
	return _SplitMain.Contract.GetETHBalance(&_SplitMain.CallOpts, account)
}

// GetETHBalance is a free data retrieval call binding the contract method 0x3bb66a7b.
//
// Solidity: function getETHBalance(address account) view returns(uint256)
func (_SplitMain *SplitMainCallerSession) GetETHBalance(account common.Address) (*big.Int, error) {
	return _SplitMain.Contract.GetETHBalance(&_SplitMain.CallOpts, account)
}

// GetHash is a free data retrieval call binding the contract method 0x1da0b8fc.
//
// Solidity: function getHash(address split) view returns(bytes32)
func (_SplitMain *SplitMainCaller) GetHash(opts *bind.CallOpts, split common.Address) ([32]byte, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "getHash", split)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetHash is a free data retrieval call binding the contract method 0x1da0b8fc.
//
// Solidity: function getHash(address split) view returns(bytes32)
func (_SplitMain *SplitMainSession) GetHash(split common.Address) ([32]byte, error) {
	return _SplitMain.Contract.GetHash(&_SplitMain.CallOpts, split)
}

// GetHash is a free data retrieval call binding the contract method 0x1da0b8fc.
//
// Solidity: function getHash(address split) view returns(bytes32)
func (_SplitMain *SplitMainCallerSession) GetHash(split common.Address) ([32]byte, error) {
	return _SplitMain.Contract.GetHash(&_SplitMain.CallOpts, split)
}

// GetNewPotentialController is a free data retrieval call binding the contract method 0xe10e51d6.
//
// Solidity: function getNewPotentialController(address split) view returns(address)
func (_SplitMain *SplitMainCaller) GetNewPotentialController(opts *bind.CallOpts, split common.Address) (common.Address, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "getNewPotentialController", split)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNewPotentialController is a free data retrieval call binding the contract method 0xe10e51d6.
//
// Solidity: function getNewPotentialController(address split) view returns(address)
func (_SplitMain *SplitMainSession) GetNewPotentialController(split common.Address) (common.Address, error) {
	return _SplitMain.Contract.GetNewPotentialController(&_SplitMain.CallOpts, split)
}

// GetNewPotentialController is a free data retrieval call binding the contract method 0xe10e51d6.
//
// Solidity: function getNewPotentialController(address split) view returns(address)
func (_SplitMain *SplitMainCallerSession) GetNewPotentialController(split common.Address) (common.Address, error) {
	return _SplitMain.Contract.GetNewPotentialController(&_SplitMain.CallOpts, split)
}

// PredictImmutableSplitAddress is a free data retrieval call binding the contract method 0x52844dd3.
//
// Solidity: function predictImmutableSplitAddress(address[] accounts, uint32[] percentAllocations, uint32 distributorFee) view returns(address split)
func (_SplitMain *SplitMainCaller) PredictImmutableSplitAddress(opts *bind.CallOpts, accounts []common.Address, percentAllocations []uint32, distributorFee uint32) (common.Address, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "predictImmutableSplitAddress", accounts, percentAllocations, distributorFee)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PredictImmutableSplitAddress is a free data retrieval call binding the contract method 0x52844dd3.
//
// Solidity: function predictImmutableSplitAddress(address[] accounts, uint32[] percentAllocations, uint32 distributorFee) view returns(address split)
func (_SplitMain *SplitMainSession) PredictImmutableSplitAddress(accounts []common.Address, percentAllocations []uint32, distributorFee uint32) (common.Address, error) {
	return _SplitMain.Contract.PredictImmutableSplitAddress(&_SplitMain.CallOpts, accounts, percentAllocations, distributorFee)
}

// PredictImmutableSplitAddress is a free data retrieval call binding the contract method 0x52844dd3.
//
// Solidity: function predictImmutableSplitAddress(address[] accounts, uint32[] percentAllocations, uint32 distributorFee) view returns(address split)
func (_SplitMain *SplitMainCallerSession) PredictImmutableSplitAddress(accounts []common.Address, percentAllocations []uint32, distributorFee uint32) (common.Address, error) {
	return _SplitMain.Contract.PredictImmutableSplitAddress(&_SplitMain.CallOpts, accounts, percentAllocations, distributorFee)
}

// WalletImplementation is a free data retrieval call binding the contract method 0x8117abc1.
//
// Solidity: function walletImplementation() view returns(address)
func (_SplitMain *SplitMainCaller) WalletImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SplitMain.contract.Call(opts, &out, "walletImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WalletImplementation is a free data retrieval call binding the contract method 0x8117abc1.
//
// Solidity: function walletImplementation() view returns(address)
func (_SplitMain *SplitMainSession) WalletImplementation() (common.Address, error) {
	return _SplitMain.Contract.WalletImplementation(&_SplitMain.CallOpts)
}

// WalletImplementation is a free data retrieval call binding the contract method 0x8117abc1.
//
// Solidity: function walletImplementation() view returns(address)
func (_SplitMain *SplitMainCallerSession) WalletImplementation() (common.Address, error) {
	return _SplitMain.Contract.WalletImplementation(&_SplitMain.CallOpts)
}

// AcceptControl is a paid mutator transaction binding the contract method 0xc7de6440.
//
// Solidity: function acceptControl(address split) returns()
func (_SplitMain *SplitMainTransactor) AcceptControl(opts *bind.TransactOpts, split common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "acceptControl", split)
}

// AcceptControl is a paid mutator transaction binding the contract method 0xc7de6440.
//
// Solidity: function acceptControl(address split) returns()
func (_SplitMain *SplitMainSession) AcceptControl(split common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.AcceptControl(&_SplitMain.TransactOpts, split)
}

// AcceptControl is a paid mutator transaction binding the contract method 0xc7de6440.
//
// Solidity: function acceptControl(address split) returns()
func (_SplitMain *SplitMainTransactorSession) AcceptControl(split common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.AcceptControl(&_SplitMain.TransactOpts, split)
}

// CancelControlTransfer is a paid mutator transaction binding the contract method 0x1267c6da.
//
// Solidity: function cancelControlTransfer(address split) returns()
func (_SplitMain *SplitMainTransactor) CancelControlTransfer(opts *bind.TransactOpts, split common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "cancelControlTransfer", split)
}

// CancelControlTransfer is a paid mutator transaction binding the contract method 0x1267c6da.
//
// Solidity: function cancelControlTransfer(address split) returns()
func (_SplitMain *SplitMainSession) CancelControlTransfer(split common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.CancelControlTransfer(&_SplitMain.TransactOpts, split)
}

// CancelControlTransfer is a paid mutator transaction binding the contract method 0x1267c6da.
//
// Solidity: function cancelControlTransfer(address split) returns()
func (_SplitMain *SplitMainTransactorSession) CancelControlTransfer(split common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.CancelControlTransfer(&_SplitMain.TransactOpts, split)
}

// CreateSplit is a paid mutator transaction binding the contract method 0x7601f782.
//
// Solidity: function createSplit(address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address controller) returns(address split)
func (_SplitMain *SplitMainTransactor) CreateSplit(opts *bind.TransactOpts, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, controller common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "createSplit", accounts, percentAllocations, distributorFee, controller)
}

// CreateSplit is a paid mutator transaction binding the contract method 0x7601f782.
//
// Solidity: function createSplit(address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address controller) returns(address split)
func (_SplitMain *SplitMainSession) CreateSplit(accounts []common.Address, percentAllocations []uint32, distributorFee uint32, controller common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.CreateSplit(&_SplitMain.TransactOpts, accounts, percentAllocations, distributorFee, controller)
}

// CreateSplit is a paid mutator transaction binding the contract method 0x7601f782.
//
// Solidity: function createSplit(address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address controller) returns(address split)
func (_SplitMain *SplitMainTransactorSession) CreateSplit(accounts []common.Address, percentAllocations []uint32, distributorFee uint32, controller common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.CreateSplit(&_SplitMain.TransactOpts, accounts, percentAllocations, distributorFee, controller)
}

// DistributeERC20 is a paid mutator transaction binding the contract method 0x15811302.
//
// Solidity: function distributeERC20(address split, address token, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactor) DistributeERC20(opts *bind.TransactOpts, split common.Address, token common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "distributeERC20", split, token, accounts, percentAllocations, distributorFee, distributorAddress)
}

// DistributeERC20 is a paid mutator transaction binding the contract method 0x15811302.
//
// Solidity: function distributeERC20(address split, address token, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainSession) DistributeERC20(split common.Address, token common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.DistributeERC20(&_SplitMain.TransactOpts, split, token, accounts, percentAllocations, distributorFee, distributorAddress)
}

// DistributeERC20 is a paid mutator transaction binding the contract method 0x15811302.
//
// Solidity: function distributeERC20(address split, address token, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactorSession) DistributeERC20(split common.Address, token common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.DistributeERC20(&_SplitMain.TransactOpts, split, token, accounts, percentAllocations, distributorFee, distributorAddress)
}

// DistributeETH is a paid mutator transaction binding the contract method 0xe61cb05e.
//
// Solidity: function distributeETH(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactor) DistributeETH(opts *bind.TransactOpts, split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "distributeETH", split, accounts, percentAllocations, distributorFee, distributorAddress)
}

// DistributeETH is a paid mutator transaction binding the contract method 0xe61cb05e.
//
// Solidity: function distributeETH(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainSession) DistributeETH(split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.DistributeETH(&_SplitMain.TransactOpts, split, accounts, percentAllocations, distributorFee, distributorAddress)
}

// DistributeETH is a paid mutator transaction binding the contract method 0xe61cb05e.
//
// Solidity: function distributeETH(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactorSession) DistributeETH(split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.DistributeETH(&_SplitMain.TransactOpts, split, accounts, percentAllocations, distributorFee, distributorAddress)
}

// MakeSplitImmutable is a paid mutator transaction binding the contract method 0x189cbaa0.
//
// Solidity: function makeSplitImmutable(address split) returns()
func (_SplitMain *SplitMainTransactor) MakeSplitImmutable(opts *bind.TransactOpts, split common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "makeSplitImmutable", split)
}

// MakeSplitImmutable is a paid mutator transaction binding the contract method 0x189cbaa0.
//
// Solidity: function makeSplitImmutable(address split) returns()
func (_SplitMain *SplitMainSession) MakeSplitImmutable(split common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.MakeSplitImmutable(&_SplitMain.TransactOpts, split)
}

// MakeSplitImmutable is a paid mutator transaction binding the contract method 0x189cbaa0.
//
// Solidity: function makeSplitImmutable(address split) returns()
func (_SplitMain *SplitMainTransactorSession) MakeSplitImmutable(split common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.MakeSplitImmutable(&_SplitMain.TransactOpts, split)
}

// TransferControl is a paid mutator transaction binding the contract method 0xd0e4b2f4.
//
// Solidity: function transferControl(address split, address newController) returns()
func (_SplitMain *SplitMainTransactor) TransferControl(opts *bind.TransactOpts, split common.Address, newController common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "transferControl", split, newController)
}

// TransferControl is a paid mutator transaction binding the contract method 0xd0e4b2f4.
//
// Solidity: function transferControl(address split, address newController) returns()
func (_SplitMain *SplitMainSession) TransferControl(split common.Address, newController common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.TransferControl(&_SplitMain.TransactOpts, split, newController)
}

// TransferControl is a paid mutator transaction binding the contract method 0xd0e4b2f4.
//
// Solidity: function transferControl(address split, address newController) returns()
func (_SplitMain *SplitMainTransactorSession) TransferControl(split common.Address, newController common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.TransferControl(&_SplitMain.TransactOpts, split, newController)
}

// UpdateAndDistributeERC20 is a paid mutator transaction binding the contract method 0x77b1e4e9.
//
// Solidity: function updateAndDistributeERC20(address split, address token, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactor) UpdateAndDistributeERC20(opts *bind.TransactOpts, split common.Address, token common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "updateAndDistributeERC20", split, token, accounts, percentAllocations, distributorFee, distributorAddress)
}

// UpdateAndDistributeERC20 is a paid mutator transaction binding the contract method 0x77b1e4e9.
//
// Solidity: function updateAndDistributeERC20(address split, address token, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainSession) UpdateAndDistributeERC20(split common.Address, token common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.UpdateAndDistributeERC20(&_SplitMain.TransactOpts, split, token, accounts, percentAllocations, distributorFee, distributorAddress)
}

// UpdateAndDistributeERC20 is a paid mutator transaction binding the contract method 0x77b1e4e9.
//
// Solidity: function updateAndDistributeERC20(address split, address token, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactorSession) UpdateAndDistributeERC20(split common.Address, token common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.UpdateAndDistributeERC20(&_SplitMain.TransactOpts, split, token, accounts, percentAllocations, distributorFee, distributorAddress)
}

// UpdateAndDistributeETH is a paid mutator transaction binding the contract method 0xa5e3909e.
//
// Solidity: function updateAndDistributeETH(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactor) UpdateAndDistributeETH(opts *bind.TransactOpts, split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "updateAndDistributeETH", split, accounts, percentAllocations, distributorFee, distributorAddress)
}

// UpdateAndDistributeETH is a paid mutator transaction binding the contract method 0xa5e3909e.
//
// Solidity: function updateAndDistributeETH(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainSession) UpdateAndDistributeETH(split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.UpdateAndDistributeETH(&_SplitMain.TransactOpts, split, accounts, percentAllocations, distributorFee, distributorAddress)
}

// UpdateAndDistributeETH is a paid mutator transaction binding the contract method 0xa5e3909e.
//
// Solidity: function updateAndDistributeETH(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee, address distributorAddress) returns()
func (_SplitMain *SplitMainTransactorSession) UpdateAndDistributeETH(split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32, distributorAddress common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.UpdateAndDistributeETH(&_SplitMain.TransactOpts, split, accounts, percentAllocations, distributorFee, distributorAddress)
}

// UpdateSplit is a paid mutator transaction binding the contract method 0xecef0ace.
//
// Solidity: function updateSplit(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee) returns()
func (_SplitMain *SplitMainTransactor) UpdateSplit(opts *bind.TransactOpts, split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "updateSplit", split, accounts, percentAllocations, distributorFee)
}

// UpdateSplit is a paid mutator transaction binding the contract method 0xecef0ace.
//
// Solidity: function updateSplit(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee) returns()
func (_SplitMain *SplitMainSession) UpdateSplit(split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32) (*types.Transaction, error) {
	return _SplitMain.Contract.UpdateSplit(&_SplitMain.TransactOpts, split, accounts, percentAllocations, distributorFee)
}

// UpdateSplit is a paid mutator transaction binding the contract method 0xecef0ace.
//
// Solidity: function updateSplit(address split, address[] accounts, uint32[] percentAllocations, uint32 distributorFee) returns()
func (_SplitMain *SplitMainTransactorSession) UpdateSplit(split common.Address, accounts []common.Address, percentAllocations []uint32, distributorFee uint32) (*types.Transaction, error) {
	return _SplitMain.Contract.UpdateSplit(&_SplitMain.TransactOpts, split, accounts, percentAllocations, distributorFee)
}

// Withdraw is a paid mutator transaction binding the contract method 0x6e5f6919.
//
// Solidity: function withdraw(address account, uint256 withdrawETH, address[] tokens) returns()
func (_SplitMain *SplitMainTransactor) Withdraw(opts *bind.TransactOpts, account common.Address, withdrawETH *big.Int, tokens []common.Address) (*types.Transaction, error) {
	return _SplitMain.contract.Transact(opts, "withdraw", account, withdrawETH, tokens)
}

// Withdraw is a paid mutator transaction binding the contract method 0x6e5f6919.
//
// Solidity: function withdraw(address account, uint256 withdrawETH, address[] tokens) returns()
func (_SplitMain *SplitMainSession) Withdraw(account common.Address, withdrawETH *big.Int, tokens []common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.Withdraw(&_SplitMain.TransactOpts, account, withdrawETH, tokens)
}

// Withdraw is a paid mutator transaction binding the contract method 0x6e5f6919.
//
// Solidity: function withdraw(address account, uint256 withdrawETH, address[] tokens) returns()
func (_SplitMain *SplitMainTransactorSession) Withdraw(account common.Address, withdrawETH *big.Int, tokens []common.Address) (*types.Transaction, error) {
	return _SplitMain.Contract.Withdraw(&_SplitMain.TransactOpts, account, withdrawETH, tokens)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SplitMain *SplitMainTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SplitMain.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SplitMain *SplitMainSession) Receive() (*types.Transaction, error) {
	return _SplitMain.Contract.Receive(&_SplitMain.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SplitMain *SplitMainTransactorSession) Receive() (*types.Transaction, error) {
	return _SplitMain.Contract.Receive(&_SplitMain.TransactOpts)
}

// SplitMainCancelControlTransferIterator is returned from FilterCancelControlTransfer and is used to iterate over the raw logs and unpacked data for CancelControlTransfer events raised by the SplitMain contract.
type SplitMainCancelControlTransferIterator struct {
	Event *SplitMainCancelControlTransfer // Event containing the contract specifics and raw log

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
func (it *SplitMainCancelControlTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainCancelControlTransfer)
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
		it.Event = new(SplitMainCancelControlTransfer)
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
func (it *SplitMainCancelControlTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainCancelControlTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainCancelControlTransfer represents a CancelControlTransfer event raised by the SplitMain contract.
type SplitMainCancelControlTransfer struct {
	Split common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCancelControlTransfer is a free log retrieval operation binding the contract event 0x6c2460a415b84be3720c209fe02f2cad7a6bcba21e8637afe8957b7ec4b6ef87.
//
// Solidity: event CancelControlTransfer(address indexed split)
func (_SplitMain *SplitMainFilterer) FilterCancelControlTransfer(opts *bind.FilterOpts, split []common.Address) (*SplitMainCancelControlTransferIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "CancelControlTransfer", splitRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainCancelControlTransferIterator{contract: _SplitMain.contract, event: "CancelControlTransfer", logs: logs, sub: sub}, nil
}

// WatchCancelControlTransfer is a free log subscription operation binding the contract event 0x6c2460a415b84be3720c209fe02f2cad7a6bcba21e8637afe8957b7ec4b6ef87.
//
// Solidity: event CancelControlTransfer(address indexed split)
func (_SplitMain *SplitMainFilterer) WatchCancelControlTransfer(opts *bind.WatchOpts, sink chan<- *SplitMainCancelControlTransfer, split []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "CancelControlTransfer", splitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainCancelControlTransfer)
				if err := _SplitMain.contract.UnpackLog(event, "CancelControlTransfer", log); err != nil {
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

// ParseCancelControlTransfer is a log parse operation binding the contract event 0x6c2460a415b84be3720c209fe02f2cad7a6bcba21e8637afe8957b7ec4b6ef87.
//
// Solidity: event CancelControlTransfer(address indexed split)
func (_SplitMain *SplitMainFilterer) ParseCancelControlTransfer(log types.Log) (*SplitMainCancelControlTransfer, error) {
	event := new(SplitMainCancelControlTransfer)
	if err := _SplitMain.contract.UnpackLog(event, "CancelControlTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainControlTransferIterator is returned from FilterControlTransfer and is used to iterate over the raw logs and unpacked data for ControlTransfer events raised by the SplitMain contract.
type SplitMainControlTransferIterator struct {
	Event *SplitMainControlTransfer // Event containing the contract specifics and raw log

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
func (it *SplitMainControlTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainControlTransfer)
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
		it.Event = new(SplitMainControlTransfer)
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
func (it *SplitMainControlTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainControlTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainControlTransfer represents a ControlTransfer event raised by the SplitMain contract.
type SplitMainControlTransfer struct {
	Split              common.Address
	PreviousController common.Address
	NewController      common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterControlTransfer is a free log retrieval operation binding the contract event 0x943d69cf2bbe08a9d44b3c4ce6da17d939d758739370620871ce99a6437866d0.
//
// Solidity: event ControlTransfer(address indexed split, address indexed previousController, address indexed newController)
func (_SplitMain *SplitMainFilterer) FilterControlTransfer(opts *bind.FilterOpts, split []common.Address, previousController []common.Address, newController []common.Address) (*SplitMainControlTransferIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}
	var previousControllerRule []interface{}
	for _, previousControllerItem := range previousController {
		previousControllerRule = append(previousControllerRule, previousControllerItem)
	}
	var newControllerRule []interface{}
	for _, newControllerItem := range newController {
		newControllerRule = append(newControllerRule, newControllerItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "ControlTransfer", splitRule, previousControllerRule, newControllerRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainControlTransferIterator{contract: _SplitMain.contract, event: "ControlTransfer", logs: logs, sub: sub}, nil
}

// WatchControlTransfer is a free log subscription operation binding the contract event 0x943d69cf2bbe08a9d44b3c4ce6da17d939d758739370620871ce99a6437866d0.
//
// Solidity: event ControlTransfer(address indexed split, address indexed previousController, address indexed newController)
func (_SplitMain *SplitMainFilterer) WatchControlTransfer(opts *bind.WatchOpts, sink chan<- *SplitMainControlTransfer, split []common.Address, previousController []common.Address, newController []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}
	var previousControllerRule []interface{}
	for _, previousControllerItem := range previousController {
		previousControllerRule = append(previousControllerRule, previousControllerItem)
	}
	var newControllerRule []interface{}
	for _, newControllerItem := range newController {
		newControllerRule = append(newControllerRule, newControllerItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "ControlTransfer", splitRule, previousControllerRule, newControllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainControlTransfer)
				if err := _SplitMain.contract.UnpackLog(event, "ControlTransfer", log); err != nil {
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

// ParseControlTransfer is a log parse operation binding the contract event 0x943d69cf2bbe08a9d44b3c4ce6da17d939d758739370620871ce99a6437866d0.
//
// Solidity: event ControlTransfer(address indexed split, address indexed previousController, address indexed newController)
func (_SplitMain *SplitMainFilterer) ParseControlTransfer(log types.Log) (*SplitMainControlTransfer, error) {
	event := new(SplitMainControlTransfer)
	if err := _SplitMain.contract.UnpackLog(event, "ControlTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainCreateSplitIterator is returned from FilterCreateSplit and is used to iterate over the raw logs and unpacked data for CreateSplit events raised by the SplitMain contract.
type SplitMainCreateSplitIterator struct {
	Event *SplitMainCreateSplit // Event containing the contract specifics and raw log

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
func (it *SplitMainCreateSplitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainCreateSplit)
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
		it.Event = new(SplitMainCreateSplit)
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
func (it *SplitMainCreateSplitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainCreateSplitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainCreateSplit represents a CreateSplit event raised by the SplitMain contract.
type SplitMainCreateSplit struct {
	Split common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCreateSplit is a free log retrieval operation binding the contract event 0x8d5f9943c664a3edaf4d3eb18cc5e2c45a7d2dc5869be33d33bbc0fff9bc2590.
//
// Solidity: event CreateSplit(address indexed split)
func (_SplitMain *SplitMainFilterer) FilterCreateSplit(opts *bind.FilterOpts, split []common.Address) (*SplitMainCreateSplitIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "CreateSplit", splitRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainCreateSplitIterator{contract: _SplitMain.contract, event: "CreateSplit", logs: logs, sub: sub}, nil
}

// WatchCreateSplit is a free log subscription operation binding the contract event 0x8d5f9943c664a3edaf4d3eb18cc5e2c45a7d2dc5869be33d33bbc0fff9bc2590.
//
// Solidity: event CreateSplit(address indexed split)
func (_SplitMain *SplitMainFilterer) WatchCreateSplit(opts *bind.WatchOpts, sink chan<- *SplitMainCreateSplit, split []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "CreateSplit", splitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainCreateSplit)
				if err := _SplitMain.contract.UnpackLog(event, "CreateSplit", log); err != nil {
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

// ParseCreateSplit is a log parse operation binding the contract event 0x8d5f9943c664a3edaf4d3eb18cc5e2c45a7d2dc5869be33d33bbc0fff9bc2590.
//
// Solidity: event CreateSplit(address indexed split)
func (_SplitMain *SplitMainFilterer) ParseCreateSplit(log types.Log) (*SplitMainCreateSplit, error) {
	event := new(SplitMainCreateSplit)
	if err := _SplitMain.contract.UnpackLog(event, "CreateSplit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainDistributeERC20Iterator is returned from FilterDistributeERC20 and is used to iterate over the raw logs and unpacked data for DistributeERC20 events raised by the SplitMain contract.
type SplitMainDistributeERC20Iterator struct {
	Event *SplitMainDistributeERC20 // Event containing the contract specifics and raw log

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
func (it *SplitMainDistributeERC20Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainDistributeERC20)
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
		it.Event = new(SplitMainDistributeERC20)
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
func (it *SplitMainDistributeERC20Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainDistributeERC20Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainDistributeERC20 represents a DistributeERC20 event raised by the SplitMain contract.
type SplitMainDistributeERC20 struct {
	Split              common.Address
	Token              common.Address
	Amount             *big.Int
	DistributorAddress common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDistributeERC20 is a free log retrieval operation binding the contract event 0xb5ee5dc3d2c31a019bbf2c787e0e9c97971c96aceea1c38c12fc8fd25c536d46.
//
// Solidity: event DistributeERC20(address indexed split, address indexed token, uint256 amount, address indexed distributorAddress)
func (_SplitMain *SplitMainFilterer) FilterDistributeERC20(opts *bind.FilterOpts, split []common.Address, token []common.Address, distributorAddress []common.Address) (*SplitMainDistributeERC20Iterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	var distributorAddressRule []interface{}
	for _, distributorAddressItem := range distributorAddress {
		distributorAddressRule = append(distributorAddressRule, distributorAddressItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "DistributeERC20", splitRule, tokenRule, distributorAddressRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainDistributeERC20Iterator{contract: _SplitMain.contract, event: "DistributeERC20", logs: logs, sub: sub}, nil
}

// WatchDistributeERC20 is a free log subscription operation binding the contract event 0xb5ee5dc3d2c31a019bbf2c787e0e9c97971c96aceea1c38c12fc8fd25c536d46.
//
// Solidity: event DistributeERC20(address indexed split, address indexed token, uint256 amount, address indexed distributorAddress)
func (_SplitMain *SplitMainFilterer) WatchDistributeERC20(opts *bind.WatchOpts, sink chan<- *SplitMainDistributeERC20, split []common.Address, token []common.Address, distributorAddress []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	var distributorAddressRule []interface{}
	for _, distributorAddressItem := range distributorAddress {
		distributorAddressRule = append(distributorAddressRule, distributorAddressItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "DistributeERC20", splitRule, tokenRule, distributorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainDistributeERC20)
				if err := _SplitMain.contract.UnpackLog(event, "DistributeERC20", log); err != nil {
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

// ParseDistributeERC20 is a log parse operation binding the contract event 0xb5ee5dc3d2c31a019bbf2c787e0e9c97971c96aceea1c38c12fc8fd25c536d46.
//
// Solidity: event DistributeERC20(address indexed split, address indexed token, uint256 amount, address indexed distributorAddress)
func (_SplitMain *SplitMainFilterer) ParseDistributeERC20(log types.Log) (*SplitMainDistributeERC20, error) {
	event := new(SplitMainDistributeERC20)
	if err := _SplitMain.contract.UnpackLog(event, "DistributeERC20", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainDistributeETHIterator is returned from FilterDistributeETH and is used to iterate over the raw logs and unpacked data for DistributeETH events raised by the SplitMain contract.
type SplitMainDistributeETHIterator struct {
	Event *SplitMainDistributeETH // Event containing the contract specifics and raw log

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
func (it *SplitMainDistributeETHIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainDistributeETH)
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
		it.Event = new(SplitMainDistributeETH)
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
func (it *SplitMainDistributeETHIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainDistributeETHIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainDistributeETH represents a DistributeETH event raised by the SplitMain contract.
type SplitMainDistributeETH struct {
	Split              common.Address
	Amount             *big.Int
	DistributorAddress common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDistributeETH is a free log retrieval operation binding the contract event 0x87c3ca0a87d9b82033e4bc55e6d30621f8d7e0c9d8ca7988edfde8932787b77b.
//
// Solidity: event DistributeETH(address indexed split, uint256 amount, address indexed distributorAddress)
func (_SplitMain *SplitMainFilterer) FilterDistributeETH(opts *bind.FilterOpts, split []common.Address, distributorAddress []common.Address) (*SplitMainDistributeETHIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	var distributorAddressRule []interface{}
	for _, distributorAddressItem := range distributorAddress {
		distributorAddressRule = append(distributorAddressRule, distributorAddressItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "DistributeETH", splitRule, distributorAddressRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainDistributeETHIterator{contract: _SplitMain.contract, event: "DistributeETH", logs: logs, sub: sub}, nil
}

// WatchDistributeETH is a free log subscription operation binding the contract event 0x87c3ca0a87d9b82033e4bc55e6d30621f8d7e0c9d8ca7988edfde8932787b77b.
//
// Solidity: event DistributeETH(address indexed split, uint256 amount, address indexed distributorAddress)
func (_SplitMain *SplitMainFilterer) WatchDistributeETH(opts *bind.WatchOpts, sink chan<- *SplitMainDistributeETH, split []common.Address, distributorAddress []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	var distributorAddressRule []interface{}
	for _, distributorAddressItem := range distributorAddress {
		distributorAddressRule = append(distributorAddressRule, distributorAddressItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "DistributeETH", splitRule, distributorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainDistributeETH)
				if err := _SplitMain.contract.UnpackLog(event, "DistributeETH", log); err != nil {
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

// ParseDistributeETH is a log parse operation binding the contract event 0x87c3ca0a87d9b82033e4bc55e6d30621f8d7e0c9d8ca7988edfde8932787b77b.
//
// Solidity: event DistributeETH(address indexed split, uint256 amount, address indexed distributorAddress)
func (_SplitMain *SplitMainFilterer) ParseDistributeETH(log types.Log) (*SplitMainDistributeETH, error) {
	event := new(SplitMainDistributeETH)
	if err := _SplitMain.contract.UnpackLog(event, "DistributeETH", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainInitiateControlTransferIterator is returned from FilterInitiateControlTransfer and is used to iterate over the raw logs and unpacked data for InitiateControlTransfer events raised by the SplitMain contract.
type SplitMainInitiateControlTransferIterator struct {
	Event *SplitMainInitiateControlTransfer // Event containing the contract specifics and raw log

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
func (it *SplitMainInitiateControlTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainInitiateControlTransfer)
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
		it.Event = new(SplitMainInitiateControlTransfer)
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
func (it *SplitMainInitiateControlTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainInitiateControlTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainInitiateControlTransfer represents a InitiateControlTransfer event raised by the SplitMain contract.
type SplitMainInitiateControlTransfer struct {
	Split                  common.Address
	NewPotentialController common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterInitiateControlTransfer is a free log retrieval operation binding the contract event 0x107cf6ea8668d533df1aab5bb8b6315bb0c25f0b6c955558d09368f290668fc7.
//
// Solidity: event InitiateControlTransfer(address indexed split, address indexed newPotentialController)
func (_SplitMain *SplitMainFilterer) FilterInitiateControlTransfer(opts *bind.FilterOpts, split []common.Address, newPotentialController []common.Address) (*SplitMainInitiateControlTransferIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}
	var newPotentialControllerRule []interface{}
	for _, newPotentialControllerItem := range newPotentialController {
		newPotentialControllerRule = append(newPotentialControllerRule, newPotentialControllerItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "InitiateControlTransfer", splitRule, newPotentialControllerRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainInitiateControlTransferIterator{contract: _SplitMain.contract, event: "InitiateControlTransfer", logs: logs, sub: sub}, nil
}

// WatchInitiateControlTransfer is a free log subscription operation binding the contract event 0x107cf6ea8668d533df1aab5bb8b6315bb0c25f0b6c955558d09368f290668fc7.
//
// Solidity: event InitiateControlTransfer(address indexed split, address indexed newPotentialController)
func (_SplitMain *SplitMainFilterer) WatchInitiateControlTransfer(opts *bind.WatchOpts, sink chan<- *SplitMainInitiateControlTransfer, split []common.Address, newPotentialController []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}
	var newPotentialControllerRule []interface{}
	for _, newPotentialControllerItem := range newPotentialController {
		newPotentialControllerRule = append(newPotentialControllerRule, newPotentialControllerItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "InitiateControlTransfer", splitRule, newPotentialControllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainInitiateControlTransfer)
				if err := _SplitMain.contract.UnpackLog(event, "InitiateControlTransfer", log); err != nil {
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

// ParseInitiateControlTransfer is a log parse operation binding the contract event 0x107cf6ea8668d533df1aab5bb8b6315bb0c25f0b6c955558d09368f290668fc7.
//
// Solidity: event InitiateControlTransfer(address indexed split, address indexed newPotentialController)
func (_SplitMain *SplitMainFilterer) ParseInitiateControlTransfer(log types.Log) (*SplitMainInitiateControlTransfer, error) {
	event := new(SplitMainInitiateControlTransfer)
	if err := _SplitMain.contract.UnpackLog(event, "InitiateControlTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainUpdateSplitIterator is returned from FilterUpdateSplit and is used to iterate over the raw logs and unpacked data for UpdateSplit events raised by the SplitMain contract.
type SplitMainUpdateSplitIterator struct {
	Event *SplitMainUpdateSplit // Event containing the contract specifics and raw log

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
func (it *SplitMainUpdateSplitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainUpdateSplit)
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
		it.Event = new(SplitMainUpdateSplit)
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
func (it *SplitMainUpdateSplitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainUpdateSplitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainUpdateSplit represents a UpdateSplit event raised by the SplitMain contract.
type SplitMainUpdateSplit struct {
	Split common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUpdateSplit is a free log retrieval operation binding the contract event 0x45e1e99513dd915ac128b94953ca64c6375717ea1894b3114db08cdca51debd2.
//
// Solidity: event UpdateSplit(address indexed split)
func (_SplitMain *SplitMainFilterer) FilterUpdateSplit(opts *bind.FilterOpts, split []common.Address) (*SplitMainUpdateSplitIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "UpdateSplit", splitRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainUpdateSplitIterator{contract: _SplitMain.contract, event: "UpdateSplit", logs: logs, sub: sub}, nil
}

// WatchUpdateSplit is a free log subscription operation binding the contract event 0x45e1e99513dd915ac128b94953ca64c6375717ea1894b3114db08cdca51debd2.
//
// Solidity: event UpdateSplit(address indexed split)
func (_SplitMain *SplitMainFilterer) WatchUpdateSplit(opts *bind.WatchOpts, sink chan<- *SplitMainUpdateSplit, split []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "UpdateSplit", splitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainUpdateSplit)
				if err := _SplitMain.contract.UnpackLog(event, "UpdateSplit", log); err != nil {
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

// ParseUpdateSplit is a log parse operation binding the contract event 0x45e1e99513dd915ac128b94953ca64c6375717ea1894b3114db08cdca51debd2.
//
// Solidity: event UpdateSplit(address indexed split)
func (_SplitMain *SplitMainFilterer) ParseUpdateSplit(log types.Log) (*SplitMainUpdateSplit, error) {
	event := new(SplitMainUpdateSplit)
	if err := _SplitMain.contract.UnpackLog(event, "UpdateSplit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SplitMainWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the SplitMain contract.
type SplitMainWithdrawalIterator struct {
	Event *SplitMainWithdrawal // Event containing the contract specifics and raw log

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
func (it *SplitMainWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitMainWithdrawal)
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
		it.Event = new(SplitMainWithdrawal)
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
func (it *SplitMainWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitMainWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitMainWithdrawal represents a Withdrawal event raised by the SplitMain contract.
type SplitMainWithdrawal struct {
	Account      common.Address
	EthAmount    *big.Int
	Tokens       []common.Address
	TokenAmounts []*big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0xa9e30bf144f83390a4fe47562a4e16892108102221c674ff538da0b72a83d174.
//
// Solidity: event Withdrawal(address indexed account, uint256 ethAmount, address[] tokens, uint256[] tokenAmounts)
func (_SplitMain *SplitMainFilterer) FilterWithdrawal(opts *bind.FilterOpts, account []common.Address) (*SplitMainWithdrawalIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SplitMain.contract.FilterLogs(opts, "Withdrawal", accountRule)
	if err != nil {
		return nil, err
	}
	return &SplitMainWithdrawalIterator{contract: _SplitMain.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0xa9e30bf144f83390a4fe47562a4e16892108102221c674ff538da0b72a83d174.
//
// Solidity: event Withdrawal(address indexed account, uint256 ethAmount, address[] tokens, uint256[] tokenAmounts)
func (_SplitMain *SplitMainFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *SplitMainWithdrawal, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _SplitMain.contract.WatchLogs(opts, "Withdrawal", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitMainWithdrawal)
				if err := _SplitMain.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0xa9e30bf144f83390a4fe47562a4e16892108102221c674ff538da0b72a83d174.
//
// Solidity: event Withdrawal(address indexed account, uint256 ethAmount, address[] tokens, uint256[] tokenAmounts)
func (_SplitMain *SplitMainFilterer) ParseWithdrawal(log types.Log) (*SplitMainWithdrawal, error) {
	event := new(SplitMainWithdrawal)
	if err := _SplitMain.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
