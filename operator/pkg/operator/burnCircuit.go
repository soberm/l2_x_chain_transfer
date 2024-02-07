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

	BatchSize = 2

	NumberBlockchains = 1
	BlockchainID      = 0
)

type BurnCircuit struct {
	Sender [BatchSize]AccountConstraints

	MerkleProofSender [BatchSize]merkle.MerkleProof

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

func (circuit *BurnCircuit) AllocateSlicesMerkleProofs() {

	for i := 0; i < BatchSize; i++ {
		circuit.MerkleProofSender[i].Path = make([]frontend.Variable, StateTreeDepth)
		//circuit.MerkleProofTransfers[i].Path = make([]frontend.Variable, TransactionsTreeDepth)
	}

}

func (circuit *BurnCircuit) Define(api frontend.API) error {

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	intermediateRoot := circuit.PreStateRoot

	leaves := make([]frontend.Variable, BatchSize)
	for i := 0; i < BatchSize; i++ {
		leaves[i] = circuit.Transfers[i].Hash(&hFunc)

		isContained := 0
		for blockchain := range circuit.Blockchains {
			if blockchain == BlockchainID {
				isContained = 1
			}
		}
		api.AssertIsEqual(isContained, 1)

		err = circuit.verifyTransferSignature(api, circuit.Transfers[i], hFunc)
		if err != nil {
			return fmt.Errorf("failed to verify transfer signature: %v", err)
		}

		intermediateRoot = circuit.burn(api, &hFunc, intermediateRoot, &circuit.Transfers[i], &circuit.Sender[i], &circuit.MerkleProofSender[i])
	}

	transactionsRoot := ComputeRoot(api, &hFunc, leaves)

	api.AssertIsEqual(transactionsRoot, circuit.TransactionsRoot)

	api.AssertIsEqual(intermediateRoot, circuit.PostStateRoot)

	return nil
}

func (circuit *BurnCircuit) verifyTransferSignature(api frontend.API, t TransferConstraints, hFunc mimc.MiMC) error {
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

func (circuit *BurnCircuit) burn(api frontend.API, hFunc hash.FieldHasher, root frontend.Variable, t *TransferConstraints, a *AccountConstraints, merkleProof *merkle.MerkleProof) frontend.Variable {
	api.Println("Burning tokens...")

	api.AssertIsEqual(merkleProof.Path[0], a.Hash(hFunc))
	api.AssertIsEqual(merkleProof.RootHash, root)

	merkleProof.VerifyProof(api, hFunc, a.Index)

	api.AssertIsEqual(a.Nonce, t.Nonce)
	api.AssertIsEqual(a.PubKey.A.X, t.SenderPubKey.A.X)
	api.AssertIsEqual(a.PubKey.A.Y, t.SenderPubKey.A.Y)

	sum := api.Add(t.Amount, t.Fee)
	api.AssertIsLessOrEqual(sum, a.Balance)

	a.Nonce = api.Add(a.Nonce, 1)
	a.Balance = api.Sub(a.Balance, sum)

	merkleProof.Path[0] = a.Hash(hFunc)

	return ComputeRootFromPath(api, merkleProof, hFunc, a.Index)
}
