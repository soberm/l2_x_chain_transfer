// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface Oracle {
    function submitTransactionsRoot(uint256 transactionsRoot, uint256 operator) external;
    function getTransactionsRoot(uint256 transactionsRoot) external view returns (uint256);
}

contract OracleMock {

    mapping(uint256 => uint256) public transactionsRoots;

    constructor(){

    }

    function submitTransactionsRoot(uint256 transactionsRoot, uint256 operator) external {
        transactionsRoots[transactionsRoot] = operator;
    }

    function getTransactionsRoot(uint256 transactionsRoot) public view returns (uint256) {
        return transactionsRoots[transactionsRoot];
    }
}
