package operator

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/hash/mimc"
)

type ClaimCircuit struct {
	Receiver       [BatchSize]AccountConstraints
	SourceOperator AccountConstraints
	TargetOperator AccountConstraints

	MerkleProofSourceOperator merkle.MerkleProof
	MerkleProofTargetOperator merkle.MerkleProof
	MerkleProofReceiver       [BatchSize]merkle.MerkleProof

	Transfers [BatchSize]TransferConstraints

	PreStateRoot     frontend.Variable `gnark:",public"`
	PostStateRoot    frontend.Variable `gnark:",public"`
	TransactionsRoot frontend.Variable `gnark:",public"`
}

func (circuit *ClaimCircuit) AllocateSlicesMerkleProofs() {

	circuit.MerkleProofSourceOperator.Path = make([]frontend.Variable, StateTreeDepth)
	circuit.MerkleProofTargetOperator.Path = make([]frontend.Variable, StateTreeDepth)

	for i := 0; i < BatchSize; i++ {
		circuit.MerkleProofReceiver[i].Path = make([]frontend.Variable, StateTreeDepth)
	}
}

func (circuit *ClaimCircuit) Define(api frontend.API) error {

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	intermediateRoot := circuit.PreStateRoot

	leaves := make([]frontend.Variable, BatchSize)
	for i := 0; i < BatchSize; i++ {
		api.AssertIsEqual(circuit.Transfers[i].Destination, BlockchainID)
		leaves[i] = circuit.Transfers[i].Hash(&hFunc)

		intermediateRoot = circuit.claim(api, &hFunc, intermediateRoot, &circuit.Transfers[i], &circuit.Receiver[i], &circuit.MerkleProofReceiver[i])
	}

	intermediateRoot = circuit.rewardOperator(api, &hFunc, intermediateRoot, circuit.Transfers, &circuit.SourceOperator, &circuit.MerkleProofSourceOperator)

	intermediateRoot = circuit.rewardOperator(api, &hFunc, intermediateRoot, circuit.Transfers, &circuit.TargetOperator, &circuit.MerkleProofTargetOperator)

	transactionsRoot := ComputeRoot(api, &hFunc, leaves)

	api.AssertIsEqual(transactionsRoot, circuit.TransactionsRoot)

	api.AssertIsEqual(intermediateRoot, circuit.PostStateRoot)

	return nil
}

func (circuit *ClaimCircuit) claim(api frontend.API, hFunc hash.FieldHasher, root frontend.Variable, t *TransferConstraints, a *AccountConstraints, merkleProof *merkle.MerkleProof) frontend.Variable {
	api.Println("Claiming tokens...")

	api.AssertIsEqual(a.PubKey.A.X, t.ReceiverPubKey.A.X)
	api.AssertIsEqual(a.PubKey.A.Y, t.ReceiverPubKey.A.Y)

	api.AssertIsEqual(merkleProof.RootHash, root)
	api.AssertIsEqual(merkleProof.Path[0], a.Hash(hFunc))
	merkleProof.VerifyProof(api, hFunc, a.Index)

	a.Balance = api.Add(a.Balance, t.Amount)
	merkleProof.Path[0] = a.Hash(hFunc)

	return ComputeRootFromPath(api, merkleProof, hFunc, a.Index)
}

func (circuit *ClaimCircuit) rewardOperator(api frontend.API, hFunc hash.FieldHasher, root frontend.Variable, t [2]TransferConstraints, a *AccountConstraints, merkleProof *merkle.MerkleProof) frontend.Variable {
	api.Println("Rewarding operator...")

	api.AssertIsEqual(merkleProof.RootHash, root)
	api.AssertIsEqual(merkleProof.Path[0], a.Hash(hFunc))
	merkleProof.VerifyProof(api, hFunc, a.Index)

	for i := 0; i < BatchSize; i++ {
		result, err := api.Compiler().NewHint(Div, 1, t[i].Fee, 2)
		if err != nil {
			return err
		}
		a.Balance = api.Add(a.Balance, result[0])
	}

	merkleProof.Path[0] = a.Hash(hFunc)

	return ComputeRootFromPath(api, merkleProof, hFunc, a.Index)
}
