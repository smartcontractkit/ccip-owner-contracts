// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gethwrappers

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

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

var CallProxyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"TargetSet\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516103e43803806103e483398181016040528101906100329190610106565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250507f3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a196816040516100959190610142565b60405180910390a15061015d565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100d3826100a8565b9050919050565b6100e3816100c8565b81146100ee57600080fd5b50565b600081519050610100816100da565b92915050565b60006020828403121561011c5761011b6100a3565b5b600061012a848285016100f1565b91505092915050565b61013c816100c8565b82525050565b60006020820190506101576000830184610133565b92915050565b60805161026d610177600039600060ac015261026d6000f3fe608060405260046000369050108061006957506336568abe60e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916600036906100479190610135565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614155b6100a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161009f90610217565b60405180910390fd5b60007f00000000000000000000000000000000000000000000000000000000000000009050366000803760008036600034855af13d6000803e80156100ec573d6000f35b3d6000fd5b600082905092915050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600082821b905092915050565b600061014183836100f1565b8261014c81356100fc565b9250600482101561018c576101877fffffffff0000000000000000000000000000000000000000000000000000000083600403600802610128565b831692505b505092915050565b600082825260208201905092915050565b7f43616c6c50726f78793a2072656e6f756e6365526f6c6520697320626c6f636b60008201527f6564000000000000000000000000000000000000000000000000000000000000602082015250565b6000610201602283610194565b915061020c826101a5565b604082019050919050565b60006020820190508181036000830152610230816101f4565b905091905056fea2646970667358221220ceb956f2ce4bedb87a33d30f718fb7372f607e76464e64446f747ffc0a5d93c064736f6c63430008130033",
}

var CallProxyABI = CallProxyMetaData.ABI

var CallProxyBin = CallProxyMetaData.Bin

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
	return address, tx, &CallProxy{address: address, abi: *parsed, CallProxyCaller: CallProxyCaller{contract: contract}, CallProxyTransactor: CallProxyTransactor{contract: contract}, CallProxyFilterer: CallProxyFilterer{contract: contract}}, nil
}

type CallProxy struct {
	address common.Address
	abi     abi.ABI
	CallProxyCaller
	CallProxyTransactor
	CallProxyFilterer
}

type CallProxyCaller struct {
	contract *bind.BoundContract
}

type CallProxyTransactor struct {
	contract *bind.BoundContract
}

type CallProxyFilterer struct {
	contract *bind.BoundContract
}

type CallProxySession struct {
	Contract     *CallProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CallProxyCallerSession struct {
	Contract *CallProxyCaller
	CallOpts bind.CallOpts
}

type CallProxyTransactorSession struct {
	Contract     *CallProxyTransactor
	TransactOpts bind.TransactOpts
}

type CallProxyRaw struct {
	Contract *CallProxy
}

type CallProxyCallerRaw struct {
	Contract *CallProxyCaller
}

type CallProxyTransactorRaw struct {
	Contract *CallProxyTransactor
}

func NewCallProxy(address common.Address, backend bind.ContractBackend) (*CallProxy, error) {
	abi, err := abi.JSON(strings.NewReader(CallProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCallProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CallProxy{address: address, abi: abi, CallProxyCaller: CallProxyCaller{contract: contract}, CallProxyTransactor: CallProxyTransactor{contract: contract}, CallProxyFilterer: CallProxyFilterer{contract: contract}}, nil
}

func NewCallProxyCaller(address common.Address, caller bind.ContractCaller) (*CallProxyCaller, error) {
	contract, err := bindCallProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CallProxyCaller{contract: contract}, nil
}

func NewCallProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*CallProxyTransactor, error) {
	contract, err := bindCallProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CallProxyTransactor{contract: contract}, nil
}

func NewCallProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*CallProxyFilterer, error) {
	contract, err := bindCallProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CallProxyFilterer{contract: contract}, nil
}

func bindCallProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CallProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CallProxy *CallProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CallProxy.Contract.CallProxyCaller.contract.Call(opts, result, method, params...)
}

func (_CallProxy *CallProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CallProxy.Contract.CallProxyTransactor.contract.Transfer(opts)
}

func (_CallProxy *CallProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CallProxy.Contract.CallProxyTransactor.contract.Transact(opts, method, params...)
}

func (_CallProxy *CallProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CallProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_CallProxy *CallProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CallProxy.Contract.contract.Transfer(opts)
}

func (_CallProxy *CallProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CallProxy.Contract.contract.Transact(opts, method, params...)
}

func (_CallProxy *CallProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _CallProxy.contract.RawTransact(opts, calldata)
}

func (_CallProxy *CallProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _CallProxy.Contract.Fallback(&_CallProxy.TransactOpts, calldata)
}

func (_CallProxy *CallProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _CallProxy.Contract.Fallback(&_CallProxy.TransactOpts, calldata)
}

type CallProxyTargetSetIterator struct {
	Event *CallProxyTargetSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CallProxyTargetSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

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

func (it *CallProxyTargetSetIterator) Error() error {
	return it.fail
}

func (it *CallProxyTargetSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CallProxyTargetSet struct {
	Target common.Address
	Raw    types.Log
}

func (_CallProxy *CallProxyFilterer) FilterTargetSet(opts *bind.FilterOpts) (*CallProxyTargetSetIterator, error) {

	logs, sub, err := _CallProxy.contract.FilterLogs(opts, "TargetSet")
	if err != nil {
		return nil, err
	}
	return &CallProxyTargetSetIterator{contract: _CallProxy.contract, event: "TargetSet", logs: logs, sub: sub}, nil
}

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

func (_CallProxy *CallProxyFilterer) ParseTargetSet(log types.Log) (*CallProxyTargetSet, error) {
	event := new(CallProxyTargetSet)
	if err := _CallProxy.contract.UnpackLog(event, "TargetSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CallProxy *CallProxy) ParseLog(log types.Log) (AbigenLog, error) {
	switch log.Topics[0] {
	case _CallProxy.abi.Events["TargetSet"].ID:
		return _CallProxy.ParseTargetSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CallProxyTargetSet) Topic() common.Hash {
	return common.HexToHash("0x3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a196")
}

func (_CallProxy *CallProxy) Address() common.Address {
	return _CallProxy.address
}

type CallProxyInterface interface {
	Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error)

	FilterTargetSet(opts *bind.FilterOpts) (*CallProxyTargetSetIterator, error)

	WatchTargetSet(opts *bind.WatchOpts, sink chan<- *CallProxyTargetSet) (event.Subscription, error)

	ParseTargetSet(log types.Log) (*CallProxyTargetSet, error)

	ParseLog(log types.Log) (AbigenLog, error)

	Address() common.Address
}
