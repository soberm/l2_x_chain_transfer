package simulator

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
	"math/rand"
	"operator/pkg/operator"
)

func generatePrivateKeys(number int) ([]*eddsa.PrivateKey, error) {
	privateKeys := make([]*eddsa.PrivateKey, number)
	for i := 0; i < number; i++ {
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r := rand.New(rand.NewSource(int64(i)))
		sk, err := eddsa.GenerateKey(r)
		if err != nil {
			return nil, fmt.Errorf("eddsa generate key: %w", err)
		}
		privateKeys[i] = sk
	}
	return privateKeys, nil
}

func createAccounts(privateKeys []*eddsa.PrivateKey) ([]*operator.Account, error) {
	accounts := make([]*operator.Account, len(privateKeys))
	for i, privateKey := range privateKeys {
		accounts[i] = &operator.Account{
			Index:     big.NewInt(int64(i)),
			Nonce:     big.NewInt(0),
			PublicKey: &privateKey.PublicKey,
			Balance:   big.NewInt(100),
		}
	}
	return accounts, nil
}
