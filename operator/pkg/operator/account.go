package operator

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"math/big"
)

const accountSize = 160

type Account struct {
	Index     *big.Int
	Nonce     *big.Int
	PublicKey *eddsa.PublicKey
	Balance   *big.Int
}

func (a *Account) Constraints() AccountConstraints {
	var pubKey eddsa2.PublicKey
	pubKey.Assign(tedwards.ID(ecc.BN254), a.PublicKey.Bytes())

	return AccountConstraints{
		Index:   a.Index,
		Nonce:   a.Nonce,
		Balance: a.Balance,
		PubKey:  pubKey,
	}
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

func (a *Account) Deserialize(data []byte) {

	a.Index = big.NewInt(0).SetBytes(data[:32])
	a.Nonce = big.NewInt(0).SetBytes(data[32:64])

	a.PublicKey = new(eddsa.PublicKey)

	a.PublicKey.A.X.SetZero()
	a.PublicKey.A.Y.SetOne()

	a.PublicKey.A.X.SetBytes(data[64:96])
	a.PublicKey.A.Y.SetBytes(data[96:128])

	a.Balance = big.NewInt(0).SetBytes(data[128:accountSize])
}
