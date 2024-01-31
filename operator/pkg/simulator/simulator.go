package simulator

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"math/big"
	"operator/pkg/operator"
)

type Simulator struct {
}

func NewSimulator() *Simulator {
	return &Simulator{}
}

func (s *Simulator) Run() error {
	log.Info("starting simulator")

	var circuit operator.Circuit
	circuit.MerkleProofSender.Path = make([]frontend.Variable, 3)

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return fmt.Errorf("compile circuit: %w", err)
	}
	log.Info("circuit compiled")

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf("setup circuit: %w", err)
	}
	log.Info("circuit setup completed")

	log.Infof("Generate private keys")
	privateKeys, err := generatePrivateKeys(4)
	if err != nil {
		return fmt.Errorf("setup circuit: %w", err)
	}

	log.Infof("Create accounts")
	accounts, err := createAccounts(privateKeys)
	if err != nil {
		return fmt.Errorf("setup circuit: %w", err)
	}

	state, err := operator.NewState(mimc.NewMiMC(), accounts)
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}

	preStateRoot, err := state.Root()
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}
	log.Infof("PreStateRoot: %v", big.NewInt(0).SetBytes(preStateRoot))

	transfer := operator.NewTransfer(10, privateKeys[0].PublicKey, privateKeys[1].PublicKey, 1)
	_, err = transfer.Sign(*privateKeys[0], mimc.NewMiMC())
	if err != nil {
		return fmt.Errorf("failed to sign transfer: %v", err)
	}

	sender, err := state.ReadAccount(0)
	if err != nil {
		return fmt.Errorf("read account: %w", err)
	}

	senderConstraints := sender.Constraints()

	root, merkleProof, err := state.MerkleProof(accounts[0].Index.Uint64())
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}

	senderMerkleProof := merkle.MerkleProof{
		RootHash: root,
		Path:     merkleProof,
	}

	err = s.UpdateState(state, transfer)
	if err != nil {
		return fmt.Errorf("failed to update state: %v", err)
	}

	postStateRoot, err := state.Root()
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}
	log.Infof("PostStateRoot: %v", big.NewInt(0).SetBytes(postStateRoot))

	transferConstraints := transfer.Constraints()

	assignment := operator.Circuit{
		Sender:            senderConstraints,
		MerkleProofSender: senderMerkleProof,
		Transfer:          transferConstraints,
		PreStateRoot:      preStateRoot,
		PostStateRoot:     postStateRoot,
	}

	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return fmt.Errorf("failed to generate proof: %v", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return fmt.Errorf("failed to verify proof: %v", err)
	}

	log.Info("stop simulator")
	return nil
}

func (s *Simulator) UpdateState(state *operator.State, t operator.Transfer) error {

	sender, err := state.ReadAccount(0)
	if err != nil {
		log.Fatalf("failed to read account: %v", err)
	}
	sender.Nonce.Add(sender.Nonce, big.NewInt(1))

	amount := big.NewInt(0)
	t.Amount.BigInt(amount)
	sender.Balance.Sub(sender.Balance, amount)

	err = state.WriteAccount(sender)
	if err != nil {
		log.Fatalf("failed to write account: %v", err)
	}

	return nil
}