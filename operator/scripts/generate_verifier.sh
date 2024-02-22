#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

WORKING_DIR="$SCRIPT_DIR/.."

cd "$WORKING_DIR" || exit 1

change_contract_name() {
    local SOURCE_FILE=$1
    local NEW_CONTRACT_NAME=$2
    local DEST_PATH=$3

    local NEW_FILE="${DEST_PATH}/${NEW_CONTRACT_NAME}.sol"

    sed -r "s/contract [a-zA-Z_][a-zA-Z0-9_]* /contract ${NEW_CONTRACT_NAME} /" "$SOURCE_FILE" > "$NEW_FILE"

    echo "Contract name changed and saved to $NEW_FILE"
}

change_contract_name "./build/burn_verifier.sol" "BurnVerifier" "../contracts/contracts"
change_contract_name "./build/claim_verifier.sol" "ClaimVerifier" "../contracts/contracts"