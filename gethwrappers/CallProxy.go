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
	Bin: "0x60a060405234801561001057600080fd5b506040516103e43803806103e483398181016040528101906100329190610106565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250507f3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a196816040516100959190610142565b60405180910390a15061015d565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100d3826100a8565b9050919050565b6100e3816100c8565b81146100ee57600080fd5b50565b600081519050610100816100da565b92915050565b60006020828403121561011c5761011b6100a3565b5b600061012a848285016100f1565b91505092915050565b61013c816100c8565b82525050565b60006020820190506101576000830184610133565b92915050565b60805161026d610177600039600060ac015261026d6000f3fe608060405260046000369050108061006957506336568abe60e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916600036906100479190610135565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614155b6100a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161009f90610217565b60405180910390fd5b60007f00000000000000000000000000000000000000000000000000000000000000009050366000803760008036600034855af13d6000803e80156100ec573d6000f35b3d6000fd5b600082905092915050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600082821b905092915050565b600061014183836100f1565b8261014c81356100fc565b9250600482101561018c576101877fffffffff0000000000000000000000000000000000000000000000000000000083600403600802610128565b831692505b505092915050565b600082825260208201905092915050565b7f43616c6c50726f78793a2072656e6f756e6365526f6c6520697320626c6f636b60008201527f6564000000000000000000000000000000000000000000000000000000000000602082015250565b6000610201602283610194565b915061020c826101a5565b604082019050919050565b60006020820190508181036000830152610230816101f4565b905091905056fea2646970667358221220ceb956f2ce4bedb87a33d30f718fb7372f607e76464e64446f747ffc0a5d93c064736f6c63430008130033",
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
