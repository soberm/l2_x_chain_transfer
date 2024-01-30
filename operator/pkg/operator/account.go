package operator

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
)

const accountSize = 160

type Account struct {
	Index     *big.Int
	Nonce     *big.Int
	PublicKey *eddsa.PublicKey
	Balance   *big.Int
}

func (a *Account) Serialize() []byte {
	var serialized []byte

	indexBytes := a.Index.Bytes()
	paddedIndex := make([]byte, 32)
	copy(paddedIndex[32-len(indexBytes):], indexBytes)
	serialized = append(serialized, paddedIndex...)

	nonceBytes := a.Nonce.Bytes()
	paddedNonce := make([]byte, 32)
	copy(paddedNonce[32-len(nonceBytes):], nonceBytes)
	serialized = append(serialized, paddedNonce...)

	publicKeyX := a.PublicKey.A.X.Bytes()
	publicKeyY := a.PublicKey.A.Y.Bytes()

	serialized = append(serialized, publicKeyX[:]...)
	serialized = append(serialized, publicKeyY[:]...)

	balanceBytes := a.Balance.Bytes()
	paddedBalance := make([]byte, 32)
	copy(paddedBalance[32-len(balanceBytes):], balanceBytes)
	serialized = append(serialized, paddedBalance...)

	return serialized
}
