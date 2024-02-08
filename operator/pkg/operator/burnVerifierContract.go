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

// BurnVerifierContractMetaData contains all meta data concerning the BurnVerifierContract contract.
var BurnVerifierContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ProofInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PublicInputNotInField\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"}],\"name\":\"compressProof\",\"outputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"compressed\",\"type\":\"uint256[4]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"compressedProof\",\"type\":\"uint256[4]\"},{\"internalType\":\"uint256[4]\",\"name\":\"input\",\"type\":\"uint256[4]\"}],\"name\":\"verifyCompressedProof\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256[4]\",\"name\":\"input\",\"type\":\"uint256[4]\"}],\"name\":\"verifyProof\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// BurnVerifierContractABI is the input ABI used to generate the binding from.
// Deprecated: Use BurnVerifierContractMetaData.ABI instead.
var BurnVerifierContractABI = BurnVerifierContractMetaData.ABI

// BurnVerifierContract is an auto generated Go binding around an Ethereum contract.
type BurnVerifierContract struct {
	BurnVerifierContractCaller     // Read-only binding to the contract
	BurnVerifierContractTransactor // Write-only binding to the contract
	BurnVerifierContractFilterer   // Log filterer for contract events
}

// BurnVerifierContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type BurnVerifierContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnVerifierContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BurnVerifierContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnVerifierContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BurnVerifierContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnVerifierContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BurnVerifierContractSession struct {
	Contract     *BurnVerifierContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BurnVerifierContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BurnVerifierContractCallerSession struct {
	Contract *BurnVerifierContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// BurnVerifierContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BurnVerifierContractTransactorSession struct {
	Contract     *BurnVerifierContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// BurnVerifierContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type BurnVerifierContractRaw struct {
	Contract *BurnVerifierContract // Generic contract binding to access the raw methods on
}

// BurnVerifierContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BurnVerifierContractCallerRaw struct {
	Contract *BurnVerifierContractCaller // Generic read-only contract binding to access the raw methods on
}

// BurnVerifierContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BurnVerifierContractTransactorRaw struct {
	Contract *BurnVerifierContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBurnVerifierContract creates a new instance of BurnVerifierContract, bound to a specific deployed contract.
func NewBurnVerifierContract(address common.Address, backend bind.ContractBackend) (*BurnVerifierContract, error) {
	contract, err := bindBurnVerifierContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnVerifierContract{BurnVerifierContractCaller: BurnVerifierContractCaller{contract: contract}, BurnVerifierContractTransactor: BurnVerifierContractTransactor{contract: contract}, BurnVerifierContractFilterer: BurnVerifierContractFilterer{contract: contract}}, nil
}

// NewBurnVerifierContractCaller creates a new read-only instance of BurnVerifierContract, bound to a specific deployed contract.
func NewBurnVerifierContractCaller(address common.Address, caller bind.ContractCaller) (*BurnVerifierContractCaller, error) {
	contract, err := bindBurnVerifierContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnVerifierContractCaller{contract: contract}, nil
}

// NewBurnVerifierContractTransactor creates a new write-only instance of BurnVerifierContract, bound to a specific deployed contract.
func NewBurnVerifierContractTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnVerifierContractTransactor, error) {
	contract, err := bindBurnVerifierContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnVerifierContractTransactor{contract: contract}, nil
}

// NewBurnVerifierContractFilterer creates a new log filterer instance of BurnVerifierContract, bound to a specific deployed contract.
func NewBurnVerifierContractFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnVerifierContractFilterer, error) {
	contract, err := bindBurnVerifierContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnVerifierContractFilterer{contract: contract}, nil
}

// bindBurnVerifierContract binds a generic wrapper to an already deployed contract.
func bindBurnVerifierContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnVerifierContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnVerifierContract *BurnVerifierContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnVerifierContract.Contract.BurnVerifierContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnVerifierContract *BurnVerifierContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnVerifierContract.Contract.BurnVerifierContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnVerifierContract *BurnVerifierContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnVerifierContract.Contract.BurnVerifierContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnVerifierContract *BurnVerifierContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnVerifierContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnVerifierContract *BurnVerifierContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnVerifierContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnVerifierContract *BurnVerifierContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnVerifierContract.Contract.contract.Transact(opts, method, params...)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_BurnVerifierContract *BurnVerifierContractCaller) CompressProof(opts *bind.CallOpts, proof [8]*big.Int) ([4]*big.Int, error) {
	var out []interface{}
	err := _BurnVerifierContract.contract.Call(opts, &out, "compressProof", proof)

	if err != nil {
		return *new([4]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([4]*big.Int)).(*[4]*big.Int)

	return out0, err

}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_BurnVerifierContract *BurnVerifierContractSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _BurnVerifierContract.Contract.CompressProof(&_BurnVerifierContract.CallOpts, proof)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_BurnVerifierContract *BurnVerifierContractCallerSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _BurnVerifierContract.Contract.CompressProof(&_BurnVerifierContract.CallOpts, proof)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xf2457c8d.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[4] input) view returns()
func (_BurnVerifierContract *BurnVerifierContractCaller) VerifyCompressedProof(opts *bind.CallOpts, compressedProof [4]*big.Int, input [4]*big.Int) error {
	var out []interface{}
	err := _BurnVerifierContract.contract.Call(opts, &out, "verifyCompressedProof", compressedProof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xf2457c8d.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[4] input) view returns()
func (_BurnVerifierContract *BurnVerifierContractSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [4]*big.Int) error {
	return _BurnVerifierContract.Contract.VerifyCompressedProof(&_BurnVerifierContract.CallOpts, compressedProof, input)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0xf2457c8d.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[4] input) view returns()
func (_BurnVerifierContract *BurnVerifierContractCallerSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [4]*big.Int) error {
	return _BurnVerifierContract.Contract.VerifyCompressedProof(&_BurnVerifierContract.CallOpts, compressedProof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x23572511.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[4] input) view returns()
func (_BurnVerifierContract *BurnVerifierContractCaller) VerifyProof(opts *bind.CallOpts, proof [8]*big.Int, input [4]*big.Int) error {
	var out []interface{}
	err := _BurnVerifierContract.contract.Call(opts, &out, "verifyProof", proof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyProof is a free data retrieval call binding the contract method 0x23572511.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[4] input) view returns()
func (_BurnVerifierContract *BurnVerifierContractSession) VerifyProof(proof [8]*big.Int, input [4]*big.Int) error {
	return _BurnVerifierContract.Contract.VerifyProof(&_BurnVerifierContract.CallOpts, proof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x23572511.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[4] input) view returns()
func (_BurnVerifierContract *BurnVerifierContractCallerSession) VerifyProof(proof [8]*big.Int, input [4]*big.Int) error {
	return _BurnVerifierContract.Contract.VerifyProof(&_BurnVerifierContract.CallOpts, proof, input)
}
