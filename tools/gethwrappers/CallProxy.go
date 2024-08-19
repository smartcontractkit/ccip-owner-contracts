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
	Bin: "0x60a060405234801561001057600080fd5b5060405161013e38038061013e83398101604081905261002f91610077565b6001600160a01b03811660808190526040519081527f3bfb4bbf112628248058745a3c57e35b13369386e474b8e56c552f3063a4a1969060200160405180910390a1506100a7565b60006020828403121561008957600080fd5b81516001600160a01b03811681146100a057600080fd5b9392505050565b608051607f6100bf600039600060060152607f6000f3fe60806040527f0000000000000000000000000000000000000000000000000000000000000000366000803760008036600034855af13d6000803e80156043573d6000f35b503d6000fdfea26469706673582212202974aca3a8ae03528c7df03132603029149d639b2cd6de0ce90e33abd7a3eb9064736f6c63430008130033",
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
