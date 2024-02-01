package simulator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
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
	circuit.AllocateSlicesMerkleProofs()

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

	var senders [operator.BatchSize]operator.AccountConstraints
	var merkleProofs [operator.BatchSize]merkle.MerkleProof
	var transfersConstraints [operator.BatchSize]operator.TransferConstraints

	hFunc := mimc.NewMiMC()
	transactionData := make([]byte, hFunc.Size()*operator.BatchSize)

	for i := 0; i < operator.BatchSize; i++ {
		sender, err := state.ReadAccount(1)
		if err != nil {
			return fmt.Errorf("read account: %w", err)
		}

		senderConstraints := sender.Constraints()
		senders[i] = senderConstraints

		transfer := operator.NewTransfer(2, 1, privateKeys[sender.Index.Uint64()].PublicKey, privateKeys[sender.Index.Uint64()].PublicKey, sender.Nonce.Uint64(), 1)
		_, msg, err := transfer.Sign(*privateKeys[sender.Index.Uint64()], mimc.NewMiMC())
		if err != nil {
			return fmt.Errorf("failed to sign transfer: %v", err)
		}

		root, merkleProof, err := state.MerkleProof(sender.Index.Uint64())
		if err != nil {
			return fmt.Errorf("create state: %w", err)
		}

		senderMerkleProof := merkle.MerkleProof{
			RootHash: root,
			Path:     merkleProof,
		}

		merkleProofs[i] = senderMerkleProof

		transfersConstraints[i] = transfer.Constraints()

		copy(transactionData[i*hFunc.Size():(i+1)*hFunc.Size()], msg)

		err = s.UpdateState(state, sender.Index.Uint64(), transfer)
		if err != nil {
			return fmt.Errorf("failed to update state: %v", err)
		}
	}

	var transactionMerkleProofs [operator.BatchSize]merkle.MerkleProof
	for i := 0; i < operator.BatchSize; i++ {
		var txBuf bytes.Buffer
		_, err = txBuf.Write(transactionData)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
		root, proof, numLeaves, _ := merkletree.BuildReaderProof(&txBuf, hFunc, hFunc.Size(), uint64(i))
		if !merkletree.VerifyProof(hFunc, root, proof, uint64(i), numLeaves) {
			return errors.New("invalid merkle proof")
		}

		var path []frontend.Variable
		for j := 0; j < len(proof); j++ {
			path = append(path, proof[j])
		}

		transactionMerkleProof := merkle.MerkleProof{
			RootHash: root,
			Path:     path,
		}

		transactionMerkleProofs[i] = transactionMerkleProof
	}

	postStateRoot, err := state.Root()
	if err != nil {
		return fmt.Errorf("create state: %w", err)
	}
	log.Infof("PostStateRoot: %v", big.NewInt(0).SetBytes(postStateRoot))

	assignment := operator.Circuit{
		Sender:               senders,
		MerkleProofSender:    merkleProofs,
		MerkleProofTransfers: transactionMerkleProofs,
		Transfers:            transfersConstraints,
		PreStateRoot:         preStateRoot,
		PostStateRoot:        postStateRoot,
		TransactionsRoot:     transactionMerkleProofs[0].RootHash,
		Blockchains:          [operator.NumberBlockchains]frontend.Variable{1},
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

func (s *Simulator) UpdateState(state *operator.State, leaf uint64, t operator.Transfer) error {

	sender, err := state.ReadAccount(leaf)
	if err != nil {
		log.Fatalf("failed to read account: %v", err)
	}
	sender.Nonce.Add(sender.Nonce, big.NewInt(1))

	amount := big.NewInt(0)
	t.Amount.BigInt(amount)

	fee := big.NewInt(0)
	t.Fee.BigInt(fee)

	sum := big.NewInt(0).Add(amount, fee)
	sender.Balance.Sub(sender.Balance, sum)

	err = state.WriteAccount(sender)
	if err != nil {
		log.Fatalf("failed to write account: %v", err)
	}

	return nil
}
