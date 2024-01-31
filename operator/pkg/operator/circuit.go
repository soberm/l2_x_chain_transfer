package operator

import (
	"fmt"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const (
	nbAccounts = 4
	depth      = 3
	BatchSize  = 2
)

type Circuit struct {
	Sender [BatchSize]AccountConstraints

	MerkleProofSender [BatchSize]merkle.MerkleProof

	Transfers [BatchSize]TransferConstraints

	PreStateRoot  frontend.Variable `gnark:",public"`
	PostStateRoot frontend.Variable `gnark:",public"`
}

type AccountConstraints struct {
	Index   frontend.Variable
	Nonce   frontend.Variable
	Balance frontend.Variable
	PubKey  eddsa.PublicKey
}

type TransferConstraints struct {
	Amount         frontend.Variable
	Nonce          frontend.Variable
	SenderPubKey   eddsa.PublicKey
	ReceiverPubKey eddsa.PublicKey
	Signature      eddsa.Signature
	//	Destination    frontend.Variable
}

func (circuit *Circuit) AllocateSlicesMerkleProofs() {

	for i := 0; i < BatchSize; i++ {
		circuit.MerkleProofSender[i].Path = make([]frontend.Variable, depth)
	}

}

func (circuit *Circuit) Define(api frontend.API) error {

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	intermediateRoot := circuit.PreStateRoot
	for i := 0; i < BatchSize; i++ {

		/*		api.Println("Sender: ", circuit.Sender[i].Index)
				api.Println("Nonce: ", circuit.Sender[i].Nonce)
				api.Println("Sender balance: ", circuit.Sender[i].Balance)
				api.Println("Sender PubKey X: ", circuit.Sender[i].PubKey.A.X)
				api.Println("Sender PubKey Y: ", circuit.Sender[i].PubKey.A.Y)*/

		hFunc.Reset()
		hFunc.Write(circuit.Sender[i].Index, circuit.Sender[i].Nonce, circuit.Sender[i].PubKey.A.X, circuit.Sender[i].PubKey.A.Y, circuit.Sender[i].Balance)

		api.AssertIsEqual(circuit.MerkleProofSender[i].Path[0], hFunc.Sum())
		api.AssertIsEqual(circuit.MerkleProofSender[i].RootHash, intermediateRoot)

		circuit.MerkleProofSender[i].VerifyProof(api, &hFunc, circuit.Sender[i].Index)

		api.AssertIsEqual(circuit.Sender[i].Nonce, circuit.Transfers[i].Nonce)
		api.AssertIsEqual(circuit.Sender[i].PubKey.A.X, circuit.Transfers[i].SenderPubKey.A.X)
		api.AssertIsEqual(circuit.Sender[i].PubKey.A.Y, circuit.Transfers[i].SenderPubKey.A.Y)

		err = verifyTransferSignature(api, circuit.Transfers[i], hFunc)
		if err != nil {
			return fmt.Errorf("failed to verify transfer signature: %v", err)
		}

		api.AssertIsLessOrEqual(circuit.Transfers[i].Amount, circuit.Sender[i].Balance)

		circuit.Sender[i].Nonce = api.Add(circuit.Sender[i].Nonce, 1)
		circuit.Sender[i].Balance = api.Sub(circuit.Sender[i].Balance, circuit.Transfers[i].Amount)

		hFunc.Reset()
		hFunc.Write(circuit.Sender[i].Index, circuit.Sender[i].Nonce, circuit.Sender[i].PubKey.A.X, circuit.Sender[i].PubKey.A.Y, circuit.Sender[i].Balance)
		circuit.MerkleProofSender[i].Path[0] = hFunc.Sum()

		intermediateRoot = ComputeRootFromPath(api, &circuit.MerkleProofSender[i], &hFunc, circuit.Sender[i].Index)
	}

	api.AssertIsEqual(intermediateRoot, circuit.PostStateRoot)

	return nil
}

func verifyTransferSignature(api frontend.API, t TransferConstraints, hFunc mimc.MiMC) error {
	hFunc.Reset()

	hFunc.Write(t.Nonce, t.Amount, t.SenderPubKey.A.X, t.SenderPubKey.A.Y, t.ReceiverPubKey.A.X, t.ReceiverPubKey.A.Y)
	htransfer := hFunc.Sum()

	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	hFunc.Reset()
	err = eddsa.Verify(curve, t.Signature, htransfer, t.SenderPubKey, &hFunc)
	if err != nil {
		return err
	}
	return nil
}
