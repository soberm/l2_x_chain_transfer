package operator

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/hash/mimc"
)

type ClaimCircuit struct {
	Receiver       [BatchSize]AccountConstraints
	SourceOperator [BatchSize]AccountConstraints
	//TargetOperator AccountConstraints

	MerkleProofSourceOperator [BatchSize]merkle.MerkleProof
	//MerkleProofTargetOperator merkle.MerkleProof
	MerkleProofReceiver  [BatchSize]merkle.MerkleProof
	MerkleProofTransfers [BatchSize]merkle.MerkleProof

	Transfers [BatchSize]TransferConstraints

	PreStateRoot     frontend.Variable `gnark:",public"`
	PostStateRoot    frontend.Variable `gnark:",public"`
	TransactionsRoot frontend.Variable `gnark:",public"`
}

func (circuit *ClaimCircuit) AllocateSlicesMerkleProofs() {

	//circuit.MerkleProofTargetOperator.Path = make([]frontend.Variable, StateTreeDepth)

	for i := 0; i < BatchSize; i++ {
		circuit.MerkleProofSourceOperator[i].Path = make([]frontend.Variable, StateTreeDepth)
		circuit.MerkleProofReceiver[i].Path = make([]frontend.Variable, StateTreeDepth)
		circuit.MerkleProofTransfers[i].Path = make([]frontend.Variable, TransactionsTreeDepth)
	}
}

func (circuit *ClaimCircuit) Define(api frontend.API) error {

	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	intermediateRoot := circuit.PreStateRoot
	for i := 0; i < BatchSize; i++ {
		api.AssertIsEqual(circuit.Transfers[i].Destination, BlockchainID)

		circuit.transferIncluded(api, &hFunc, &circuit.Transfers[i], &circuit.Receiver[i], &circuit.MerkleProofTransfers[i], i)

		/*		api.Println("Executing transfer ", i, "...")
				api.Println("Sender: ", circuit.Receiver[i].Index)
				api.Println("Nonce: ", circuit.Receiver[i].Nonce)
				api.Println("Sender balance: ", circuit.Receiver[i].Balance)
				api.Println("Sender PubKey X: ", circuit.Receiver[i].PubKey.A.X)
				api.Println("Sender PubKey Y: ", circuit.Receiver[i].PubKey.A.Y)*/

		intermediateRoot = circuit.claim(api, &hFunc, intermediateRoot, &circuit.Transfers[i], &circuit.Receiver[i], &circuit.MerkleProofReceiver[i])

		intermediateRoot = circuit.rewardOperator(api, &hFunc, intermediateRoot, &circuit.Transfers[i], &circuit.SourceOperator[i], &circuit.MerkleProofSourceOperator[i])

		//intermediateRoot = circuit.rewardOperator(api, &hFunc, intermediateRoot, &circuit.Transfers[i], &circuit.TargetOperator, &circuit.MerkleProofTargetOperator)
	}
	api.Println("Intermediate root: ", intermediateRoot)
	api.Println("Post state root: ", circuit.PostStateRoot)
	api.AssertIsEqual(intermediateRoot, circuit.PostStateRoot)

	return nil
}

func (circuit *ClaimCircuit) transferIncluded(api frontend.API, hFunc hash.FieldHasher, t *TransferConstraints, a *AccountConstraints, merkleProof *merkle.MerkleProof, index int) {
	api.AssertIsEqual(a.PubKey.A.X, t.SenderPubKey.A.X)
	api.AssertIsEqual(a.PubKey.A.Y, t.SenderPubKey.A.Y)

	api.AssertIsEqual(merkleProof.Path[0], t.Hash(hFunc))
	api.AssertIsEqual(merkleProof.RootHash, circuit.TransactionsRoot)

	merkleProof.VerifyProof(api, hFunc, index)
}

func (circuit *ClaimCircuit) claim(api frontend.API, hFunc hash.FieldHasher, root frontend.Variable, t *TransferConstraints, a *AccountConstraints, merkleProof *merkle.MerkleProof) frontend.Variable {
	api.Println("Claiming tokens...")

	api.AssertIsEqual(merkleProof.RootHash, root)
	api.AssertIsEqual(merkleProof.Path[0], a.Hash(hFunc))
	merkleProof.VerifyProof(api, hFunc, a.Index)

	a.Balance = api.Add(a.Balance, t.Amount)
	merkleProof.Path[0] = a.Hash(hFunc)

	return ComputeRootFromPath(api, merkleProof, hFunc, a.Index)
}

func (circuit *ClaimCircuit) rewardOperator(api frontend.API, hFunc hash.FieldHasher, root frontend.Variable, t *TransferConstraints, a *AccountConstraints, merkleProof *merkle.MerkleProof) frontend.Variable {
	api.Println("Rewarding operator...")

	result, err := api.Compiler().NewHint(Div, 1, t.Fee, 2)
	if err != nil {
		return err
	}

	api.AssertIsEqual(merkleProof.RootHash, root)
	api.AssertIsEqual(merkleProof.Path[0], a.Hash(hFunc))
	merkleProof.VerifyProof(api, hFunc, a.Index)

	a.Balance = api.Add(a.Balance, result[0])
	merkleProof.Path[0] = a.Hash(hFunc)

	return ComputeRootFromPath(api, merkleProof, hFunc, a.Index)
}
