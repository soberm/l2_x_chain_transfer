# Efficient Cross-Blockchain Token Transfers with Rollback Support

This project contains the source code for the paper "Efficient 
Cross-Blockchain Token Transfers with Rollback Support". We provide a
prototypiycal implementation of the smart contracts and an operator 
which simulates the execution of test transactions.

## Prerequisites

You need to have the following software installed:

* [Golang](https://golang.org/doc/install) (version 1.21.6)
* [Node.js](https://nodejs.org/) (version >= 21.6.1)
* [Ganache](https://www.trufflesuite.com/ganache) (version >= 2.5.4)
* [Solidity](https://docs.soliditylang.org/en/latest/installing-solidity.html) (^0.8.0)

## Installation

### Constraint System Setup

1. Change into the operator directory: `cd operator/`
2. Install all dependencies: `go mod download`
3. Adapt the batch size and the number of accounts in const.go
4. Build the constraint system setup: `go build -o constraint_system_setup`
5. Run the constraint system setup: `./constraint_system_setup -b ./build`
6. Generate the verifier contracts: `./scripts/generate_verifier.sh`

### Smart Contracts

1. Change into the contract directory: `cd contracts/`
2. Install all dependencies: `npm install`
3. Compile contracts: `hardhat compile`
4. Deploy contracts: `hardhat run --network <your_network> ./scripts/deploy.ts`

### Operator

1. Change into the operator directory: `cd operator/`
2. Install all dependencies: `go mod download`
3. Build the operator: `go build -o operator`
4. Run the operator: `./operator -c ./configs/config.json`

## Contributing

This is a research prototype. We welcome anyone to contribute. File a bug report or submit feature requests through the issue tracker. If you want to contribute feel free to submit a pull request.

## Acknowledgement

The financial support by the Austrian Federal Ministry for Digital and Economic Affairs, the National Foundation for Research, Technology and Development as well as the Christian Doppler Research Association is gratefully acknowledged.

## Licence

This project is licensed under the [MIT License](LICENSE).