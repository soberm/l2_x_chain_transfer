// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./BurnVerifier.sol";
import "./ClaimVerifier.sol";
import "./OracleMock.sol";

contract Rollup {

    uint256 public constant BATCH_SIZE = 128;
    uint256 public constant SUPPORTED_BLOCKCHAIN = 1;

    struct Transfer {
        uint256 nonce;
        uint256 amount;
        uint256[2] sender;
        uint256[2] receiver;
        uint256 fee;
        uint256 dest;
    }

    event BurnEvent(
        uint256 preStateRoot,
        uint256 postStateRoot,
        uint256 transactionsRoot,
        uint256[4] compressedProof
    );
    uint256 public stateRoot;

    BurnVerifier private burnVerifier;
    ClaimVerifier private claimVerifier;
    Oracle private oracle;

    mapping(uint256 => bool) private usedTransactionsRoots;

    constructor(uint256 _stateRoot, address _burnVerifier, address _claimVerifier, address _oracleAddress) {
        stateRoot = _stateRoot;
        burnVerifier = BurnVerifier(_burnVerifier);
        claimVerifier = ClaimVerifier(_claimVerifier);
        oracle = Oracle(_oracleAddress);
    }

    function Burn(
        uint256 postStateRoot,
        uint256 transactionsRoot,
        uint256[4] memory compressedProof,
        Transfer[BATCH_SIZE] calldata transfers
    ) public {
        uint[4] memory input = [stateRoot, postStateRoot, transactionsRoot, SUPPORTED_BLOCKCHAIN];
        burnVerifier.verifyCompressedProof(compressedProof, input);

        stateRoot = postStateRoot;
    }

    function Claim(
        uint256 postStateRoot,
        uint256 _transactionsRoot,
        uint256 _operator,
        uint256[4] memory compressedProof,
        Transfer[BATCH_SIZE] calldata transfers
    ) public {
        require(!usedTransactionsRoots[_transactionsRoot], "transactionsRoot already used");

        uint256 operator = oracle.getTransactionsRoot(_transactionsRoot);
        require(operator == _operator, "invalid operator");

        uint[4] memory input = [stateRoot, postStateRoot, _transactionsRoot, _operator];
        claimVerifier.verifyCompressedProof(compressedProof, input);

        stateRoot = postStateRoot;
        usedTransactionsRoots[_transactionsRoot] = true;
    }

}
