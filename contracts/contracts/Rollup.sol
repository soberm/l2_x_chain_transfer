// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./BurnVerifier.sol";
import "./ClaimVerifier.sol";

contract Rollup {

    uint256 private stateRoot;

    BurnVerifier private burnVerifier;
    ClaimVerifier private claimVerifier;

    constructor(uint256 _stateRoot, address _burnVerifier, address _claimVerifier) {
        stateRoot = _stateRoot;
        burnVerifier = BurnVerifier(_burnVerifier);
        claimVerifier = ClaimVerifier(_claimVerifier);
    }

    function Burn(
        uint256 postStateRoot,
        uint256 transactionsRoot,
        uint256[4] memory compressedProof
    ) public {
        uint[3] memory input = [stateRoot, postStateRoot, transactionsRoot];
        claimVerifier.verifyCompressedProof(compressedProof, input);

        stateRoot = postStateRoot;
    }
}
