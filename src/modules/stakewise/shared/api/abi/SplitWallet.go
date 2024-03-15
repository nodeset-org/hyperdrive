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

// SplitWalletMetaData contains all meta data concerning the SplitWallet contract.
var SplitWalletMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"split\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ReceiveETH\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendERC20ToMain\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendETHToMain\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"splitMain\",\"outputs\":[{\"internalType\":\"contractISplitMain\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SplitWalletABI is the input ABI used to generate the binding from.
// Deprecated: Use SplitWalletMetaData.ABI instead.
var SplitWalletABI = SplitWalletMetaData.ABI

// SplitWallet is an auto generated Go binding around an Ethereum contract.
type SplitWallet struct {
	SplitWalletCaller     // Read-only binding to the contract
	SplitWalletTransactor // Write-only binding to the contract
	SplitWalletFilterer   // Log filterer for contract events
}

// SplitWalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type SplitWalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SplitWalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SplitWalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SplitWalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SplitWalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SplitWalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SplitWalletSession struct {
	Contract     *SplitWallet      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SplitWalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SplitWalletCallerSession struct {
	Contract *SplitWalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SplitWalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SplitWalletTransactorSession struct {
	Contract     *SplitWalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SplitWalletRaw is an auto generated low-level Go binding around an Ethereum contract.
type SplitWalletRaw struct {
	Contract *SplitWallet // Generic contract binding to access the raw methods on
}

// SplitWalletCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SplitWalletCallerRaw struct {
	Contract *SplitWalletCaller // Generic read-only contract binding to access the raw methods on
}

// SplitWalletTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SplitWalletTransactorRaw struct {
	Contract *SplitWalletTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSplitWallet creates a new instance of SplitWallet, bound to a specific deployed contract.
func NewSplitWallet(address common.Address, backend bind.ContractBackend) (*SplitWallet, error) {
	contract, err := bindSplitWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SplitWallet{SplitWalletCaller: SplitWalletCaller{contract: contract}, SplitWalletTransactor: SplitWalletTransactor{contract: contract}, SplitWalletFilterer: SplitWalletFilterer{contract: contract}}, nil
}

// NewSplitWalletCaller creates a new read-only instance of SplitWallet, bound to a specific deployed contract.
func NewSplitWalletCaller(address common.Address, caller bind.ContractCaller) (*SplitWalletCaller, error) {
	contract, err := bindSplitWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SplitWalletCaller{contract: contract}, nil
}

// NewSplitWalletTransactor creates a new write-only instance of SplitWallet, bound to a specific deployed contract.
func NewSplitWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*SplitWalletTransactor, error) {
	contract, err := bindSplitWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SplitWalletTransactor{contract: contract}, nil
}

// NewSplitWalletFilterer creates a new log filterer instance of SplitWallet, bound to a specific deployed contract.
func NewSplitWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*SplitWalletFilterer, error) {
	contract, err := bindSplitWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SplitWalletFilterer{contract: contract}, nil
}

// bindSplitWallet binds a generic wrapper to an already deployed contract.
func bindSplitWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SplitWalletMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SplitWallet *SplitWalletRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SplitWallet.Contract.SplitWalletCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SplitWallet *SplitWalletRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SplitWallet.Contract.SplitWalletTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SplitWallet *SplitWalletRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SplitWallet.Contract.SplitWalletTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SplitWallet *SplitWalletCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SplitWallet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SplitWallet *SplitWalletTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SplitWallet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SplitWallet *SplitWalletTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SplitWallet.Contract.contract.Transact(opts, method, params...)
}

// SplitMain is a free data retrieval call binding the contract method 0x0e769b2b.
//
// Solidity: function splitMain() view returns(address)
func (_SplitWallet *SplitWalletCaller) SplitMain(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SplitWallet.contract.Call(opts, &out, "splitMain")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SplitMain is a free data retrieval call binding the contract method 0x0e769b2b.
//
// Solidity: function splitMain() view returns(address)
func (_SplitWallet *SplitWalletSession) SplitMain() (common.Address, error) {
	return _SplitWallet.Contract.SplitMain(&_SplitWallet.CallOpts)
}

// SplitMain is a free data retrieval call binding the contract method 0x0e769b2b.
//
// Solidity: function splitMain() view returns(address)
func (_SplitWallet *SplitWalletCallerSession) SplitMain() (common.Address, error) {
	return _SplitWallet.Contract.SplitMain(&_SplitWallet.CallOpts)
}

// SendERC20ToMain is a paid mutator transaction binding the contract method 0x7c1f3ffe.
//
// Solidity: function sendERC20ToMain(address token, uint256 amount) payable returns()
func (_SplitWallet *SplitWalletTransactor) SendERC20ToMain(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SplitWallet.contract.Transact(opts, "sendERC20ToMain", token, amount)
}

// SendERC20ToMain is a paid mutator transaction binding the contract method 0x7c1f3ffe.
//
// Solidity: function sendERC20ToMain(address token, uint256 amount) payable returns()
func (_SplitWallet *SplitWalletSession) SendERC20ToMain(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SplitWallet.Contract.SendERC20ToMain(&_SplitWallet.TransactOpts, token, amount)
}

// SendERC20ToMain is a paid mutator transaction binding the contract method 0x7c1f3ffe.
//
// Solidity: function sendERC20ToMain(address token, uint256 amount) payable returns()
func (_SplitWallet *SplitWalletTransactorSession) SendERC20ToMain(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SplitWallet.Contract.SendERC20ToMain(&_SplitWallet.TransactOpts, token, amount)
}

// SendETHToMain is a paid mutator transaction binding the contract method 0xab0ebff4.
//
// Solidity: function sendETHToMain(uint256 amount) payable returns()
func (_SplitWallet *SplitWalletTransactor) SendETHToMain(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _SplitWallet.contract.Transact(opts, "sendETHToMain", amount)
}

// SendETHToMain is a paid mutator transaction binding the contract method 0xab0ebff4.
//
// Solidity: function sendETHToMain(uint256 amount) payable returns()
func (_SplitWallet *SplitWalletSession) SendETHToMain(amount *big.Int) (*types.Transaction, error) {
	return _SplitWallet.Contract.SendETHToMain(&_SplitWallet.TransactOpts, amount)
}

// SendETHToMain is a paid mutator transaction binding the contract method 0xab0ebff4.
//
// Solidity: function sendETHToMain(uint256 amount) payable returns()
func (_SplitWallet *SplitWalletTransactorSession) SendETHToMain(amount *big.Int) (*types.Transaction, error) {
	return _SplitWallet.Contract.SendETHToMain(&_SplitWallet.TransactOpts, amount)
}

// SplitWalletReceiveETHIterator is returned from FilterReceiveETH and is used to iterate over the raw logs and unpacked data for ReceiveETH events raised by the SplitWallet contract.
type SplitWalletReceiveETHIterator struct {
	Event *SplitWalletReceiveETH // Event containing the contract specifics and raw log

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
func (it *SplitWalletReceiveETHIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SplitWalletReceiveETH)
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
		it.Event = new(SplitWalletReceiveETH)
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
func (it *SplitWalletReceiveETHIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SplitWalletReceiveETHIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SplitWalletReceiveETH represents a ReceiveETH event raised by the SplitWallet contract.
type SplitWalletReceiveETH struct {
	Split  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReceiveETH is a free log retrieval operation binding the contract event 0x830d2d700a97af574b186c80d40429385d24241565b08a7c559ba283a964d9b1.
//
// Solidity: event ReceiveETH(address indexed split, uint256 amount)
func (_SplitWallet *SplitWalletFilterer) FilterReceiveETH(opts *bind.FilterOpts, split []common.Address) (*SplitWalletReceiveETHIterator, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitWallet.contract.FilterLogs(opts, "ReceiveETH", splitRule)
	if err != nil {
		return nil, err
	}
	return &SplitWalletReceiveETHIterator{contract: _SplitWallet.contract, event: "ReceiveETH", logs: logs, sub: sub}, nil
}

// WatchReceiveETH is a free log subscription operation binding the contract event 0x830d2d700a97af574b186c80d40429385d24241565b08a7c559ba283a964d9b1.
//
// Solidity: event ReceiveETH(address indexed split, uint256 amount)
func (_SplitWallet *SplitWalletFilterer) WatchReceiveETH(opts *bind.WatchOpts, sink chan<- *SplitWalletReceiveETH, split []common.Address) (event.Subscription, error) {

	var splitRule []interface{}
	for _, splitItem := range split {
		splitRule = append(splitRule, splitItem)
	}

	logs, sub, err := _SplitWallet.contract.WatchLogs(opts, "ReceiveETH", splitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SplitWalletReceiveETH)
				if err := _SplitWallet.contract.UnpackLog(event, "ReceiveETH", log); err != nil {
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

// ParseReceiveETH is a log parse operation binding the contract event 0x830d2d700a97af574b186c80d40429385d24241565b08a7c559ba283a964d9b1.
//
// Solidity: event ReceiveETH(address indexed split, uint256 amount)
func (_SplitWallet *SplitWalletFilterer) ParseReceiveETH(log types.Log) (*SplitWalletReceiveETH, error) {
	event := new(SplitWalletReceiveETH)
	if err := _SplitWallet.contract.UnpackLog(event, "ReceiveETH", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
