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

// ClaimVerifierContractMetaData contains all meta data concerning the ClaimVerifierContract contract.
var ClaimVerifierContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ProofInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PublicInputNotInField\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"}],\"name\":\"compressProof\",\"outputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"compressed\",\"type\":\"uint256[4]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"compressedProof\",\"type\":\"uint256[4]\"},{\"internalType\":\"uint256[3]\",\"name\":\"input\",\"type\":\"uint256[3]\"}],\"name\":\"verifyCompressedProof\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256[3]\",\"name\":\"input\",\"type\":\"uint256[3]\"}],\"name\":\"verifyProof\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ClaimVerifierContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ClaimVerifierContractMetaData.ABI instead.
var ClaimVerifierContractABI = ClaimVerifierContractMetaData.ABI

// ClaimVerifierContract is an auto generated Go binding around an Ethereum contract.
type ClaimVerifierContract struct {
	ClaimVerifierContractCaller     // Read-only binding to the contract
	ClaimVerifierContractTransactor // Write-only binding to the contract
	ClaimVerifierContractFilterer   // Log filterer for contract events
}

// ClaimVerifierContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClaimVerifierContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimVerifierContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClaimVerifierContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimVerifierContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClaimVerifierContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimVerifierContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClaimVerifierContractSession struct {
	Contract     *ClaimVerifierContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ClaimVerifierContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClaimVerifierContractCallerSession struct {
	Contract *ClaimVerifierContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ClaimVerifierContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClaimVerifierContractTransactorSession struct {
	Contract     *ClaimVerifierContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ClaimVerifierContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClaimVerifierContractRaw struct {
	Contract *ClaimVerifierContract // Generic contract binding to access the raw methods on
}

// ClaimVerifierContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClaimVerifierContractCallerRaw struct {
	Contract *ClaimVerifierContractCaller // Generic read-only contract binding to access the raw methods on
}

// ClaimVerifierContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClaimVerifierContractTransactorRaw struct {
	Contract *ClaimVerifierContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClaimVerifierContract creates a new instance of ClaimVerifierContract, bound to a specific deployed contract.
func NewClaimVerifierContract(address common.Address, backend bind.ContractBackend) (*ClaimVerifierContract, error) {
	contract, err := bindClaimVerifierContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ClaimVerifierContract{ClaimVerifierContractCaller: ClaimVerifierContractCaller{contract: contract}, ClaimVerifierContractTransactor: ClaimVerifierContractTransactor{contract: contract}, ClaimVerifierContractFilterer: ClaimVerifierContractFilterer{contract: contract}}, nil
}

// NewClaimVerifierContractCaller creates a new read-only instance of ClaimVerifierContract, bound to a specific deployed contract.
func NewClaimVerifierContractCaller(address common.Address, caller bind.ContractCaller) (*ClaimVerifierContractCaller, error) {
	contract, err := bindClaimVerifierContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClaimVerifierContractCaller{contract: contract}, nil
}

// NewClaimVerifierContractTransactor creates a new write-only instance of ClaimVerifierContract, bound to a specific deployed contract.
func NewClaimVerifierContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ClaimVerifierContractTransactor, error) {
	contract, err := bindClaimVerifierContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClaimVerifierContractTransactor{contract: contract}, nil
}

// NewClaimVerifierContractFilterer creates a new log filterer instance of ClaimVerifierContract, bound to a specific deployed contract.
func NewClaimVerifierContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ClaimVerifierContractFilterer, error) {
	contract, err := bindClaimVerifierContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClaimVerifierContractFilterer{contract: contract}, nil
}

// bindClaimVerifierContract binds a generic wrapper to an already deployed contract.
func bindClaimVerifierContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ClaimVerifierContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClaimVerifierContract *ClaimVerifierContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClaimVerifierContract.Contract.ClaimVerifierContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClaimVerifierContract *ClaimVerifierContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimVerifierContract.Contract.ClaimVerifierContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClaimVerifierContract *ClaimVerifierContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClaimVerifierContract.Contract.ClaimVerifierContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClaimVerifierContract *ClaimVerifierContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClaimVerifierContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClaimVerifierContract *ClaimVerifierContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimVerifierContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClaimVerifierContract *ClaimVerifierContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClaimVerifierContract.Contract.contract.Transact(opts, method, params...)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_ClaimVerifierContract *ClaimVerifierContractCaller) CompressProof(opts *bind.CallOpts, proof [8]*big.Int) ([4]*big.Int, error) {
	var out []interface{}
	err := _ClaimVerifierContract.contract.Call(opts, &out, "compressProof", proof)

	if err != nil {
		return *new([4]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([4]*big.Int)).(*[4]*big.Int)

	return out0, err

}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_ClaimVerifierContract *ClaimVerifierContractSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _ClaimVerifierContract.Contract.CompressProof(&_ClaimVerifierContract.CallOpts, proof)
}

// CompressProof is a free data retrieval call binding the contract method 0x44f63692.
//
// Solidity: function compressProof(uint256[8] proof) view returns(uint256[4] compressed)
func (_ClaimVerifierContract *ClaimVerifierContractCallerSession) CompressProof(proof [8]*big.Int) ([4]*big.Int, error) {
	return _ClaimVerifierContract.Contract.CompressProof(&_ClaimVerifierContract.CallOpts, proof)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0x98c13db9.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[3] input) view returns()
func (_ClaimVerifierContract *ClaimVerifierContractCaller) VerifyCompressedProof(opts *bind.CallOpts, compressedProof [4]*big.Int, input [3]*big.Int) error {
	var out []interface{}
	err := _ClaimVerifierContract.contract.Call(opts, &out, "verifyCompressedProof", compressedProof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0x98c13db9.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[3] input) view returns()
func (_ClaimVerifierContract *ClaimVerifierContractSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [3]*big.Int) error {
	return _ClaimVerifierContract.Contract.VerifyCompressedProof(&_ClaimVerifierContract.CallOpts, compressedProof, input)
}

// VerifyCompressedProof is a free data retrieval call binding the contract method 0x98c13db9.
//
// Solidity: function verifyCompressedProof(uint256[4] compressedProof, uint256[3] input) view returns()
func (_ClaimVerifierContract *ClaimVerifierContractCallerSession) VerifyCompressedProof(compressedProof [4]*big.Int, input [3]*big.Int) error {
	return _ClaimVerifierContract.Contract.VerifyCompressedProof(&_ClaimVerifierContract.CallOpts, compressedProof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x65c03259.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[3] input) view returns()
func (_ClaimVerifierContract *ClaimVerifierContractCaller) VerifyProof(opts *bind.CallOpts, proof [8]*big.Int, input [3]*big.Int) error {
	var out []interface{}
	err := _ClaimVerifierContract.contract.Call(opts, &out, "verifyProof", proof, input)

	if err != nil {
		return err
	}

	return err

}

// VerifyProof is a free data retrieval call binding the contract method 0x65c03259.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[3] input) view returns()
func (_ClaimVerifierContract *ClaimVerifierContractSession) VerifyProof(proof [8]*big.Int, input [3]*big.Int) error {
	return _ClaimVerifierContract.Contract.VerifyProof(&_ClaimVerifierContract.CallOpts, proof, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x65c03259.
//
// Solidity: function verifyProof(uint256[8] proof, uint256[3] input) view returns()
func (_ClaimVerifierContract *ClaimVerifierContractCallerSession) VerifyProof(proof [8]*big.Int, input [3]*big.Int) error {
	return _ClaimVerifierContract.Contract.VerifyProof(&_ClaimVerifierContract.CallOpts, proof, input)
}
