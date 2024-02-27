#!/bin/sh

cd "$(dirname "$0")" || exit 1

if ! which abigen >/dev/null; then
  echo "error: abigen not installed" >&2
  exit 1
fi

abigen --abi ../../contracts/abi/contracts/Rollup.sol/Rollup.json --pkg operator --type RollupContract --out ../pkg/operator/rollupContract.go
abigen --abi ../../contracts/abi/contracts/OracleMock.sol/OracleMock.json --pkg operator --type OracleMockContract --out ../pkg/operator/oracleMockContract.go
abigen --abi ../../contracts/abi/contracts/BurnVerifier.sol/BurnVerifier.json --pkg operator --type BurnVerifierContract --out ../pkg/operator/burnVerifierContract.go
abigen --abi ../../contracts/abi/contracts/ClaimVerifier.sol/ClaimVerifier.json --pkg operator --type ClaimVerifierContract --out ../pkg/operator/claimVerifierContract.go