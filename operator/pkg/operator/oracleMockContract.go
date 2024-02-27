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

// OracleMockContractMetaData contains all meta data concerning the OracleMockContract contract.
var OracleMockContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transactionsRoot\",\"type\":\"uint256\"}],\"name\":\"getTransactionsRoot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"transactionsRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"operator\",\"type\":\"uint256\"}],\"name\":\"submitTransactionsRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transactionsRoots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OracleMockContractABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleMockContractMetaData.ABI instead.
var OracleMockContractABI = OracleMockContractMetaData.ABI

// OracleMockContract is an auto generated Go binding around an Ethereum contract.
type OracleMockContract struct {
	OracleMockContractCaller     // Read-only binding to the contract
	OracleMockContractTransactor // Write-only binding to the contract
	OracleMockContractFilterer   // Log filterer for contract events
}

// OracleMockContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleMockContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleMockContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleMockContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleMockContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleMockContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleMockContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleMockContractSession struct {
	Contract     *OracleMockContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OracleMockContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleMockContractCallerSession struct {
	Contract *OracleMockContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// OracleMockContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleMockContractTransactorSession struct {
	Contract     *OracleMockContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OracleMockContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleMockContractRaw struct {
	Contract *OracleMockContract // Generic contract binding to access the raw methods on
}

// OracleMockContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleMockContractCallerRaw struct {
	Contract *OracleMockContractCaller // Generic read-only contract binding to access the raw methods on
}

// OracleMockContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleMockContractTransactorRaw struct {
	Contract *OracleMockContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleMockContract creates a new instance of OracleMockContract, bound to a specific deployed contract.
func NewOracleMockContract(address common.Address, backend bind.ContractBackend) (*OracleMockContract, error) {
	contract, err := bindOracleMockContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleMockContract{OracleMockContractCaller: OracleMockContractCaller{contract: contract}, OracleMockContractTransactor: OracleMockContractTransactor{contract: contract}, OracleMockContractFilterer: OracleMockContractFilterer{contract: contract}}, nil
}

// NewOracleMockContractCaller creates a new read-only instance of OracleMockContract, bound to a specific deployed contract.
func NewOracleMockContractCaller(address common.Address, caller bind.ContractCaller) (*OracleMockContractCaller, error) {
	contract, err := bindOracleMockContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleMockContractCaller{contract: contract}, nil
}

// NewOracleMockContractTransactor creates a new write-only instance of OracleMockContract, bound to a specific deployed contract.
func NewOracleMockContractTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleMockContractTransactor, error) {
	contract, err := bindOracleMockContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleMockContractTransactor{contract: contract}, nil
}

// NewOracleMockContractFilterer creates a new log filterer instance of OracleMockContract, bound to a specific deployed contract.
func NewOracleMockContractFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleMockContractFilterer, error) {
	contract, err := bindOracleMockContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleMockContractFilterer{contract: contract}, nil
}

// bindOracleMockContract binds a generic wrapper to an already deployed contract.
func bindOracleMockContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleMockContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleMockContract *OracleMockContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleMockContract.Contract.OracleMockContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleMockContract *OracleMockContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleMockContract.Contract.OracleMockContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleMockContract *OracleMockContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleMockContract.Contract.OracleMockContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleMockContract *OracleMockContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleMockContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleMockContract *OracleMockContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleMockContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleMockContract *OracleMockContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleMockContract.Contract.contract.Transact(opts, method, params...)
}

// GetTransactionsRoot is a free data retrieval call binding the contract method 0x2199a00d.
//
// Solidity: function getTransactionsRoot(uint256 transactionsRoot) view returns(uint256)
func (_OracleMockContract *OracleMockContractCaller) GetTransactionsRoot(opts *bind.CallOpts, transactionsRoot *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OracleMockContract.contract.Call(opts, &out, "getTransactionsRoot", transactionsRoot)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTransactionsRoot is a free data retrieval call binding the contract method 0x2199a00d.
//
// Solidity: function getTransactionsRoot(uint256 transactionsRoot) view returns(uint256)
func (_OracleMockContract *OracleMockContractSession) GetTransactionsRoot(transactionsRoot *big.Int) (*big.Int, error) {
	return _OracleMockContract.Contract.GetTransactionsRoot(&_OracleMockContract.CallOpts, transactionsRoot)
}

// GetTransactionsRoot is a free data retrieval call binding the contract method 0x2199a00d.
//
// Solidity: function getTransactionsRoot(uint256 transactionsRoot) view returns(uint256)
func (_OracleMockContract *OracleMockContractCallerSession) GetTransactionsRoot(transactionsRoot *big.Int) (*big.Int, error) {
	return _OracleMockContract.Contract.GetTransactionsRoot(&_OracleMockContract.CallOpts, transactionsRoot)
}

// TransactionsRoots is a free data retrieval call binding the contract method 0x77c12252.
//
// Solidity: function transactionsRoots(uint256 ) view returns(uint256)
func (_OracleMockContract *OracleMockContractCaller) TransactionsRoots(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OracleMockContract.contract.Call(opts, &out, "transactionsRoots", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransactionsRoots is a free data retrieval call binding the contract method 0x77c12252.
//
// Solidity: function transactionsRoots(uint256 ) view returns(uint256)
func (_OracleMockContract *OracleMockContractSession) TransactionsRoots(arg0 *big.Int) (*big.Int, error) {
	return _OracleMockContract.Contract.TransactionsRoots(&_OracleMockContract.CallOpts, arg0)
}

// TransactionsRoots is a free data retrieval call binding the contract method 0x77c12252.
//
// Solidity: function transactionsRoots(uint256 ) view returns(uint256)
func (_OracleMockContract *OracleMockContractCallerSession) TransactionsRoots(arg0 *big.Int) (*big.Int, error) {
	return _OracleMockContract.Contract.TransactionsRoots(&_OracleMockContract.CallOpts, arg0)
}

// SubmitTransactionsRoot is a paid mutator transaction binding the contract method 0xec898e4d.
//
// Solidity: function submitTransactionsRoot(uint256 transactionsRoot, uint256 operator) returns()
func (_OracleMockContract *OracleMockContractTransactor) SubmitTransactionsRoot(opts *bind.TransactOpts, transactionsRoot *big.Int, operator *big.Int) (*types.Transaction, error) {
	return _OracleMockContract.contract.Transact(opts, "submitTransactionsRoot", transactionsRoot, operator)
}

// SubmitTransactionsRoot is a paid mutator transaction binding the contract method 0xec898e4d.
//
// Solidity: function submitTransactionsRoot(uint256 transactionsRoot, uint256 operator) returns()
func (_OracleMockContract *OracleMockContractSession) SubmitTransactionsRoot(transactionsRoot *big.Int, operator *big.Int) (*types.Transaction, error) {
	return _OracleMockContract.Contract.SubmitTransactionsRoot(&_OracleMockContract.TransactOpts, transactionsRoot, operator)
}

// SubmitTransactionsRoot is a paid mutator transaction binding the contract method 0xec898e4d.
//
// Solidity: function submitTransactionsRoot(uint256 transactionsRoot, uint256 operator) returns()
func (_OracleMockContract *OracleMockContractTransactorSession) SubmitTransactionsRoot(transactionsRoot *big.Int, operator *big.Int) (*types.Transaction, error) {
	return _OracleMockContract.Contract.SubmitTransactionsRoot(&_OracleMockContract.TransactOpts, transactionsRoot, operator)
}
