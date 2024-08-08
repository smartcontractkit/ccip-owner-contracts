// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gethwrappers

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

// CallProxyMetaData contains all meta data concerning the CallProxy contract.
var CallProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"TargetSet\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161013e38038061013e83398101604081905261002f91610077565b6001600160a01b03811660808190526040519081527f3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a1969060200160405180910390a1506100a7565b60006020828403121561008957600080fd5b81516001600160a01b03811681146100a057600080fd5b9392505050565b608051607f6100bf600039600060060152607f6000f3fe60806040527f0000000000000000000000000000000000000000000000000000000000000000366000803760008036600034855af13d6000803e80156043573d6000f35b503d6000fdfea26469706673582212202974aca3a8ae03528c7df03132603029149d639b2cd6de0ce90e33abd7a3eb9064736f6c63430008130033",
}

// CallProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use CallProxyMetaData.ABI instead.
var CallProxyABI = CallProxyMetaData.ABI

// CallProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CallProxyMetaData.Bin instead.
var CallProxyBin = CallProxyMetaData.Bin

// DeployCallProxy deploys a new Ethereum contract, binding an instance of CallProxy to it.
func DeployCallProxy(auth *bind.TransactOpts, backend bind.ContractBackend, target common.Address) (common.Address, *types.Transaction, *CallProxy, error) {
	parsed, err := CallProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CallProxyBin), backend, target)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CallProxy{CallProxyCaller: CallProxyCaller{contract: contract}, CallProxyTransactor: CallProxyTransactor{contract: contract}, CallProxyFilterer: CallProxyFilterer{contract: contract}}, nil
}

// CallProxy is an auto generated Go binding around an Ethereum contract.
type CallProxy struct {
	CallProxyCaller     // Read-only binding to the contract
	CallProxyTransactor // Write-only binding to the contract
	CallProxyFilterer   // Log filterer for contract events
}

// CallProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type CallProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CallProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CallProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CallProxySession struct {
	Contract     *CallProxy        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CallProxyCallerSession struct {
	Contract *CallProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// CallProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CallProxyTransactorSession struct {
	Contract     *CallProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CallProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type CallProxyRaw struct {
	Contract *CallProxy // Generic contract binding to access the raw methods on
}

// CallProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CallProxyCallerRaw struct {
	Contract *CallProxyCaller // Generic read-only contract binding to access the raw methods on
}

// CallProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CallProxyTransactorRaw struct {
	Contract *CallProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCallProxy creates a new instance of CallProxy, bound to a specific deployed contract.
func NewCallProxy(address common.Address, backend bind.ContractBackend) (*CallProxy, error) {
	contract, err := bindCallProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CallProxy{CallProxyCaller: CallProxyCaller{contract: contract}, CallProxyTransactor: CallProxyTransactor{contract: contract}, CallProxyFilterer: CallProxyFilterer{contract: contract}}, nil
}

// NewCallProxyCaller creates a new read-only instance of CallProxy, bound to a specific deployed contract.
func NewCallProxyCaller(address common.Address, caller bind.ContractCaller) (*CallProxyCaller, error) {
	contract, err := bindCallProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CallProxyCaller{contract: contract}, nil
}

// NewCallProxyTransactor creates a new write-only instance of CallProxy, bound to a specific deployed contract.
func NewCallProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*CallProxyTransactor, error) {
	contract, err := bindCallProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CallProxyTransactor{contract: contract}, nil
}

// NewCallProxyFilterer creates a new log filterer instance of CallProxy, bound to a specific deployed contract.
func NewCallProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*CallProxyFilterer, error) {
	contract, err := bindCallProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CallProxyFilterer{contract: contract}, nil
}

// bindCallProxy binds a generic wrapper to an already deployed contract.
func bindCallProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CallProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CallProxy *CallProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CallProxy.Contract.CallProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CallProxy *CallProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CallProxy.Contract.CallProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CallProxy *CallProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CallProxy.Contract.CallProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CallProxy *CallProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CallProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CallProxy *CallProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CallProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CallProxy *CallProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CallProxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_CallProxy *CallProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _CallProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_CallProxy *CallProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _CallProxy.Contract.Fallback(&_CallProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_CallProxy *CallProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _CallProxy.Contract.Fallback(&_CallProxy.TransactOpts, calldata)
}

// CallProxyTargetSetIterator is returned from FilterTargetSet and is used to iterate over the raw logs and unpacked data for TargetSet events raised by the CallProxy contract.
type CallProxyTargetSetIterator struct {
	Event *CallProxyTargetSet // Event containing the contract specifics and raw log

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
func (it *CallProxyTargetSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CallProxyTargetSet)
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
		it.Event = new(CallProxyTargetSet)
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
func (it *CallProxyTargetSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CallProxyTargetSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CallProxyTargetSet represents a TargetSet event raised by the CallProxy contract.
type CallProxyTargetSet struct {
	Target common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTargetSet is a free log retrieval operation binding the contract event 0x3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a196.
//
// Solidity: event TargetSet(address target)
func (_CallProxy *CallProxyFilterer) FilterTargetSet(opts *bind.FilterOpts) (*CallProxyTargetSetIterator, error) {

	logs, sub, err := _CallProxy.contract.FilterLogs(opts, "TargetSet")
	if err != nil {
		return nil, err
	}
	return &CallProxyTargetSetIterator{contract: _CallProxy.contract, event: "TargetSet", logs: logs, sub: sub}, nil
}

// WatchTargetSet is a free log subscription operation binding the contract event 0x3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a196.
//
// Solidity: event TargetSet(address target)
func (_CallProxy *CallProxyFilterer) WatchTargetSet(opts *bind.WatchOpts, sink chan<- *CallProxyTargetSet) (event.Subscription, error) {

	logs, sub, err := _CallProxy.contract.WatchLogs(opts, "TargetSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CallProxyTargetSet)
				if err := _CallProxy.contract.UnpackLog(event, "TargetSet", log); err != nil {
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

// ParseTargetSet is a log parse operation binding the contract event 0x3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a196.
//
// Solidity: event TargetSet(address target)
func (_CallProxy *CallProxyFilterer) ParseTargetSet(log types.Log) (*CallProxyTargetSet, error) {
	event := new(CallProxyTargetSet)
	if err := _CallProxy.contract.UnpackLog(event, "TargetSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
