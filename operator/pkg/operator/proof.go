package operator

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark/backend/groth16"
	"math/big"
)

func ProofToEthereumProof(p groth16.Proof) ([8]*big.Int, error) {

	var proof [8]*big.Int

	var buf bytes.Buffer
	_, err := p.WriteRawTo(&buf)
	if err != nil {
		return proof, fmt.Errorf("write raw proof to: %w", err)
	}
	proofBytes := buf.Bytes()

	for i := 0; i < len(proof); i++ {
		proof[i] = new(big.Int).SetBytes(proofBytes[fp.Bytes*i : fp.Bytes*(i+1)])
	}

	return proof, nil
}
