package operator

import (
	"fmt"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const (
	NumberAccounts = 4
	StateTreeDepth = 3

	BatchSize             = 2
	TransactionsTreeDepth = 2

	NumberBlockchains = 1
	BlockchainID      = 0
)

type Circuit struct {
	Sender [BatchSize]AccountConstraints

	MerkleProofSender    [BatchSize]merkle.MerkleProof
	MerkleProofTransfers [BatchSize]merkle.MerkleProof

	Transfers [BatchSize]TransferConstraints

	PreStateRoot     frontend.Variable                    `gnark:",public"`
	PostStateRoot    frontend.Variable                    `gnark:",public"`
	TransactionsRoot frontend.Variable                    `gnark:",public"`
	Blockchains      [NumberBlockchains]frontend.Variable `gnark:",public"`
}

type AccountConstraints struct {
	Index   frontend.Variable
	Nonce   frontend.Variable
	Balance frontend.Variable
	PubKey  eddsa.PublicKey
}

type TransferConstraints struct {
	Amount         frontend.Variable
	Fee            frontend.Variable
	Nonce          frontend.Variable
	SenderPubKey   eddsa.PublicKey
	ReceiverPubKey eddsa.PublicKey
	Signature      eddsa.Signature
	Destination    frontend.Variable
}

func (a *AccountConstraints) Hash(h hash.FieldHasher) frontend.Variable {
	h.Reset()
	h.Write(a.Index, a.Nonce, a.PubKey.A.X, a.PubKey.A.Y, a.Balance)
	return h.Sum()
}

func (t *TransferConstraints) Hash(h hash.FieldHasher) frontend.Variable {
	h.Reset()
	h.Write(t.Nonce, t.Amount, t.Fee, t.SenderPubKey.A.X, t.SenderPubKey.A.Y, t.ReceiverPubKey.A.X, t.ReceiverPubKey.A.Y, t.Destination)
	return h.Sum()
}

func (circuit *Circuit) AllocateSlicesMerkleProofs() {

	for i := 0; i < BatchSize; i++ {
		circuit.MerkleProofSender[i].Path = make([]frontend.Variable, StateTreeDepth)
		circuit.MerkleProofTransfers[i].Path = make([]frontend.Variable, TransactionsTreeDepth)
	}

}

func (circuit *Circuit) Define(api frontend.API) error {

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	intermediateRoot := circuit.PreStateRoot
	for i := 0; i < BatchSize; i++ {
		api.Println("Executing transfer ", i, "...")
		/*		api.Println("Sender: ", circuit.Sender[i].Index)
				api.Println("Nonce: ", circuit.Sender[i].Nonce)
				api.Println("Sender balance: ", circuit.Sender[i].Balance)
				api.Println("Sender PubKey X: ", circuit.Sender[i].PubKey.A.X)
				api.Println("Sender PubKey Y: ", circuit.Sender[i].PubKey.A.Y)*/

		isContained := 0
		for blockchain := range circuit.Blockchains {
			if blockchain == BlockchainID {
				isContained = 1
			}
		}
		api.AssertIsEqual(isContained, 1)

		api.AssertIsEqual(circuit.MerkleProofSender[i].Path[0], circuit.Sender[i].Hash(&hFunc))
		api.AssertIsEqual(circuit.MerkleProofSender[i].RootHash, intermediateRoot)

		circuit.MerkleProofSender[i].VerifyProof(api, &hFunc, circuit.Sender[i].Index)

		api.AssertIsEqual(circuit.Sender[i].Nonce, circuit.Transfers[i].Nonce)
		api.AssertIsEqual(circuit.Sender[i].PubKey.A.X, circuit.Transfers[i].SenderPubKey.A.X)
		api.AssertIsEqual(circuit.Sender[i].PubKey.A.Y, circuit.Transfers[i].SenderPubKey.A.Y)

		api.AssertIsEqual(circuit.MerkleProofTransfers[i].Path[0], circuit.Transfers[i].Hash(&hFunc))
		api.AssertIsEqual(circuit.MerkleProofTransfers[i].RootHash, circuit.TransactionsRoot)

		circuit.MerkleProofTransfers[i].VerifyProof(api, &hFunc, i)

		err = circuit.verifyTransferSignature(api, circuit.Transfers[i], hFunc)
		if err != nil {
			return fmt.Errorf("failed to verify transfer signature: %v", err)
		}

		circuit.burn(api, &circuit.Transfers[i], &circuit.Sender[i])

		circuit.MerkleProofSender[i].Path[0] = circuit.Sender[i].Hash(&hFunc)

		intermediateRoot = ComputeRootFromPath(api, &circuit.MerkleProofSender[i], &hFunc, circuit.Sender[i].Index)
	}

	api.AssertIsEqual(intermediateRoot, circuit.PostStateRoot)

	return nil
}

func (circuit *Circuit) verifyTransferSignature(api frontend.API, t TransferConstraints, hFunc mimc.MiMC) error {
	api.Println("Verifying signature...")
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return fmt.Errorf("failed to create twisted edwards curve: %v", err)
	}

	tHash := t.Hash(&hFunc)
	hFunc.Reset()

	err = eddsa.Verify(curve, t.Signature, tHash, t.SenderPubKey, &hFunc)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}
	return nil
}

func (circuit *Circuit) burn(api frontend.API, t *TransferConstraints, a *AccountConstraints) {
	api.Println("Burning tokens...")
	sum := api.Add(t.Amount, t.Fee)
	api.AssertIsLessOrEqual(sum, a.Balance)

	a.Nonce = api.Add(a.Nonce, 1)
	a.Balance = api.Sub(a.Balance, sum)
}
