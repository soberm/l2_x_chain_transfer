// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package operator

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

// RollupContractMetaData contains all meta data concerning the RollupContract contract.
var RollupContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_stateRoot\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_burnVerifier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_claimVerifier\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preStateRoot\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"postStateRoot\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transactionsRoot\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[4]\",\"name\":\"compressedProof\",\"type\":\"uint256[4]\"}],\"name\":\"BurnEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"postStateRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionsRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256[4]\",\"name\":\"compressedProof\",\"type\":\"uint256[4]\"}],\"name\":\"Burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"postStateRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionsRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256[4]\",\"name\":\"compressedProof\",\"type\":\"uint256[4]\"}],\"name\":\"Claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stateRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// RollupContractABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupContractMetaData.ABI instead.
var RollupContractABI = RollupContractMetaData.ABI

// RollupContract is an auto generated Go binding around an Ethereum contract.
type RollupContract struct {
	RollupContractCaller     // Read-only binding to the contract
	RollupContractTransactor // Write-only binding to the contract
	RollupContractFilterer   // Log filterer for contract events
}

// RollupContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupContractSession struct {
	Contract     *RollupContract   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupContractCallerSession struct {
	Contract *RollupContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RollupContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupContractTransactorSession struct {
	Contract     *RollupContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RollupContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupContractRaw struct {
	Contract *RollupContract // Generic contract binding to access the raw methods on
}

// RollupContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupContractCallerRaw struct {
	Contract *RollupContractCaller // Generic read-only contract binding to access the raw methods on
}

// RollupContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupContractTransactorRaw struct {
	Contract *RollupContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollupContract creates a new instance of RollupContract, bound to a specific deployed contract.
func NewRollupContract(address common.Address, backend bind.ContractBackend) (*RollupContract, error) {
	contract, err := bindRollupContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollupContract{RollupContractCaller: RollupContractCaller{contract: contract}, RollupContractTransactor: RollupContractTransactor{contract: contract}, RollupContractFilterer: RollupContractFilterer{contract: contract}}, nil
}

// NewRollupContractCaller creates a new read-only instance of RollupContract, bound to a specific deployed contract.
func NewRollupContractCaller(address common.Address, caller bind.ContractCaller) (*RollupContractCaller, error) {
	contract, err := bindRollupContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupContractCaller{contract: contract}, nil
}

// NewRollupContractTransactor creates a new write-only instance of RollupContract, bound to a specific deployed contract.
func NewRollupContractTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupContractTransactor, error) {
	contract, err := bindRollupContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupContractTransactor{contract: contract}, nil
}

// NewRollupContractFilterer creates a new log filterer instance of RollupContract, bound to a specific deployed contract.
func NewRollupContractFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupContractFilterer, error) {
	contract, err := bindRollupContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupContractFilterer{contract: contract}, nil
}

// bindRollupContract binds a generic wrapper to an already deployed contract.
func bindRollupContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RollupContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupContract *RollupContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupContract.Contract.RollupContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupContract *RollupContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupContract.Contract.RollupContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupContract *RollupContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupContract.Contract.RollupContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupContract *RollupContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupContract *RollupContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupContract *RollupContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupContract.Contract.contract.Transact(opts, method, params...)
}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(uint256)
func (_RollupContract *RollupContractCaller) StateRoot(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupContract.contract.Call(opts, &out, "stateRoot")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(uint256)
func (_RollupContract *RollupContractSession) StateRoot() (*big.Int, error) {
	return _RollupContract.Contract.StateRoot(&_RollupContract.CallOpts)
}

// StateRoot is a free data retrieval call binding the contract method 0x9588eca2.
//
// Solidity: function stateRoot() view returns(uint256)
func (_RollupContract *RollupContractCallerSession) StateRoot() (*big.Int, error) {
	return _RollupContract.Contract.StateRoot(&_RollupContract.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x150d98ea.
//
// Solidity: function Burn(uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof) returns()
func (_RollupContract *RollupContractTransactor) Burn(opts *bind.TransactOpts, postStateRoot *big.Int, transactionsRoot *big.Int, compressedProof [4]*big.Int) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "Burn", postStateRoot, transactionsRoot, compressedProof)
}

// Burn is a paid mutator transaction binding the contract method 0x150d98ea.
//
// Solidity: function Burn(uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof) returns()
func (_RollupContract *RollupContractSession) Burn(postStateRoot *big.Int, transactionsRoot *big.Int, compressedProof [4]*big.Int) (*types.Transaction, error) {
	return _RollupContract.Contract.Burn(&_RollupContract.TransactOpts, postStateRoot, transactionsRoot, compressedProof)
}

// Burn is a paid mutator transaction binding the contract method 0x150d98ea.
//
// Solidity: function Burn(uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof) returns()
func (_RollupContract *RollupContractTransactorSession) Burn(postStateRoot *big.Int, transactionsRoot *big.Int, compressedProof [4]*big.Int) (*types.Transaction, error) {
	return _RollupContract.Contract.Burn(&_RollupContract.TransactOpts, postStateRoot, transactionsRoot, compressedProof)
}

// Claim is a paid mutator transaction binding the contract method 0x9a7602bd.
//
// Solidity: function Claim(uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof) returns()
func (_RollupContract *RollupContractTransactor) Claim(opts *bind.TransactOpts, postStateRoot *big.Int, transactionsRoot *big.Int, compressedProof [4]*big.Int) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "Claim", postStateRoot, transactionsRoot, compressedProof)
}

// Claim is a paid mutator transaction binding the contract method 0x9a7602bd.
//
// Solidity: function Claim(uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof) returns()
func (_RollupContract *RollupContractSession) Claim(postStateRoot *big.Int, transactionsRoot *big.Int, compressedProof [4]*big.Int) (*types.Transaction, error) {
	return _RollupContract.Contract.Claim(&_RollupContract.TransactOpts, postStateRoot, transactionsRoot, compressedProof)
}

// Claim is a paid mutator transaction binding the contract method 0x9a7602bd.
//
// Solidity: function Claim(uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof) returns()
func (_RollupContract *RollupContractTransactorSession) Claim(postStateRoot *big.Int, transactionsRoot *big.Int, compressedProof [4]*big.Int) (*types.Transaction, error) {
	return _RollupContract.Contract.Claim(&_RollupContract.TransactOpts, postStateRoot, transactionsRoot, compressedProof)
}

// RollupContractBurnEventIterator is returned from FilterBurnEvent and is used to iterate over the raw logs and unpacked data for BurnEvent events raised by the RollupContract contract.
type RollupContractBurnEventIterator struct {
	Event *RollupContractBurnEvent // Event containing the contract specifics and raw log

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
func (it *RollupContractBurnEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupContractBurnEvent)
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
		it.Event = new(RollupContractBurnEvent)
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
func (it *RollupContractBurnEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupContractBurnEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupContractBurnEvent represents a BurnEvent event raised by the RollupContract contract.
type RollupContractBurnEvent struct {
	PreStateRoot     *big.Int
	PostStateRoot    *big.Int
	TransactionsRoot *big.Int
	CompressedProof  [4]*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBurnEvent is a free log retrieval operation binding the contract event 0xa80ab886d69a255698da9bf9787ab6837dc01363efa653efbffbb9fa69058734.
//
// Solidity: event BurnEvent(uint256 preStateRoot, uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof)
func (_RollupContract *RollupContractFilterer) FilterBurnEvent(opts *bind.FilterOpts) (*RollupContractBurnEventIterator, error) {

	logs, sub, err := _RollupContract.contract.FilterLogs(opts, "BurnEvent")
	if err != nil {
		return nil, err
	}
	return &RollupContractBurnEventIterator{contract: _RollupContract.contract, event: "BurnEvent", logs: logs, sub: sub}, nil
}

// WatchBurnEvent is a free log subscription operation binding the contract event 0xa80ab886d69a255698da9bf9787ab6837dc01363efa653efbffbb9fa69058734.
//
// Solidity: event BurnEvent(uint256 preStateRoot, uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof)
func (_RollupContract *RollupContractFilterer) WatchBurnEvent(opts *bind.WatchOpts, sink chan<- *RollupContractBurnEvent) (event.Subscription, error) {

	logs, sub, err := _RollupContract.contract.WatchLogs(opts, "BurnEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupContractBurnEvent)
				if err := _RollupContract.contract.UnpackLog(event, "BurnEvent", log); err != nil {
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

// ParseBurnEvent is a log parse operation binding the contract event 0xa80ab886d69a255698da9bf9787ab6837dc01363efa653efbffbb9fa69058734.
//
// Solidity: event BurnEvent(uint256 preStateRoot, uint256 postStateRoot, uint256 transactionsRoot, uint256[4] compressedProof)
func (_RollupContract *RollupContractFilterer) ParseBurnEvent(log types.Log) (*RollupContractBurnEvent, error) {
	event := new(RollupContractBurnEvent)
	if err := _RollupContract.contract.UnpackLog(event, "BurnEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
