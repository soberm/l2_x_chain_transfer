package operator

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"math/big"
	"math/rand"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

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

func createAccounts(privateKeys []*eddsa.PrivateKey) ([]*Account, error) {
	accounts := make([]*Account, len(privateKeys))
	for i, privateKey := range privateKeys {
		accounts[i] = &Account{
			Index:     big.NewInt(int64(i)),
			Nonce:     big.NewInt(0),
			PublicKey: &privateKey.PublicKey,
			Balance:   big.NewInt(1000000),
		}
	}
	return accounts, nil
}
