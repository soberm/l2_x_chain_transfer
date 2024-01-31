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
	nbAccounts       = 4
	depth            = 3
	BatchSizeCircuit = 1
)

type Circuit struct {
	Sender AccountConstraints

	MerkleProofSender merkle.MerkleProof

	Transfer TransferConstraints

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

func (circuit *Circuit) Define(api frontend.API) error {

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	api.AssertIsEqual(circuit.MerkleProofSender.RootHash, circuit.PreStateRoot)

	circuit.MerkleProofSender.VerifyProof(api, &hFunc, circuit.Sender.Index)
	err = verifyTransferSignature(api, circuit.Transfer, hFunc)
	if err != nil {
		return fmt.Errorf("failed to verify transfer signature: %v", err)
	}

	circuit.Sender.Nonce = api.Add(circuit.Sender.Nonce, 1)
	api.AssertIsEqual(circuit.Sender.Nonce, circuit.Transfer.Nonce)

	api.AssertIsLessOrEqual(circuit.Transfer.Amount, circuit.Sender.Balance)

	circuit.Sender.Balance = api.Sub(circuit.Sender.Balance, circuit.Transfer.Amount)

	hFunc.Reset()
	hFunc.Write(circuit.Sender.Index, circuit.Sender.Nonce, circuit.Sender.PubKey.A.X, circuit.Sender.PubKey.A.Y, circuit.Sender.Balance)
	circuit.MerkleProofSender.Path[0] = hFunc.Sum()

	intermediateRoot := ComputeRootFromPath(api, &circuit.MerkleProofSender, &hFunc, circuit.Sender.Index)
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
