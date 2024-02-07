package simulator

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"math/big"
	"operator/pkg/operator"
)

type Rollup struct {
	Accounts    []*operator.Account
	PrivateKeys []*eddsa.PrivateKey
	State       *operator.State
}

func NewRollup() (*Rollup, error) {

	privateKeys, err := generatePrivateKeys(operator.NumberAccounts)
	if err != nil {
		return nil, fmt.Errorf("generate private keys: %w", err)
	}

	log.Infof("Create accounts")
	accounts, err := createAccounts(privateKeys)
	if err != nil {
		return nil, fmt.Errorf("create accounts: %w", err)
	}

	state, err := operator.NewState(mimc.NewMiMC(), accounts)
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	return &Rollup{Accounts: accounts, PrivateKeys: privateKeys, State: state}, nil
}

func (r *Rollup) GenerateTransfers(number int) ([]operator.Transfer, error) {
	hFunc := mimc.NewMiMC()
	transfers := make([]operator.Transfer, number)
	transferData := make([]byte, hFunc.Size()*number)

	for i := 0; i < number; i++ {
		sender, err := r.State.ReadAccount(uint64(i))
		if err != nil {
			return nil, fmt.Errorf("read account: %w", err)
		}
		transfer := operator.NewTransfer(10,
			4,
			r.PrivateKeys[sender.Index.Uint64()].PublicKey,
			r.PrivateKeys[sender.Index.Uint64()].PublicKey,
			sender.Nonce.Uint64(),
			0,
		)

		_, msg, err := transfer.Sign(*r.PrivateKeys[sender.Index.Uint64()], mimc.NewMiMC())
		if err != nil {
			return nil, fmt.Errorf("failed to sign transfer: %v", err)
		}

		transfers[i] = transfer
		copy(transferData[i*hFunc.Size():(i+1)*hFunc.Size()], msg)
	}

	return transfers, nil
}

func (r *Rollup) Burn(transfers []operator.Transfer) (witness.Witness, error) {

	preStateRoot, err := r.State.Root()
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	var senderConstraints [operator.BatchSize]operator.AccountConstraints
	var senderMerkleProofs [operator.BatchSize]merkle.MerkleProof
	var transfersConstraints [operator.BatchSize]operator.TransferConstraints

	transfersRoot, err := r.TransfersRoot(transfers)
	if err != nil {
		return nil, fmt.Errorf("generate transactions: %w", err)
	}

	for i := 0; i < len(transfers); i++ {
		sender, err := r.State.ReadAccount(uint64(i))
		if err != nil {
			return nil, fmt.Errorf("read account: %w", err)
		}

		senderConstraints[i] = sender.Constraints()

		root, senderMerkleProof, err := r.State.MerkleProof(sender.Index.Uint64())
		if err != nil {
			return nil, fmt.Errorf("create state: %w", err)
		}

		senderMerkleProofs[i] = merkle.MerkleProof{
			RootHash: root,
			Path:     senderMerkleProof,
		}

		transfersConstraints[i] = transfers[i].Constraints()

		err = r.UpdateSender(&sender, &transfers[i])
		if err != nil {
			return nil, fmt.Errorf("failed to write account: %v", err)
		}
	}

	postStateRoot, err := r.State.Root()
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	assignment := &operator.BurnCircuit{
		Sender:            senderConstraints,
		MerkleProofSender: senderMerkleProofs,
		Transfers:         transfersConstraints,
		PreStateRoot:      preStateRoot,
		PostStateRoot:     postStateRoot,
		TransactionsRoot:  transfersRoot,
		Blockchains:       [operator.NumberBlockchains]frontend.Variable{1},
	}

	w, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to create witness: %v", err)
	}

	return w, nil
}

func (r *Rollup) Claim(transfers []operator.Transfer) (witness.Witness, error) {
	preStateRoot, err := r.State.Root()
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	var receiverConstraints [operator.BatchSize]operator.AccountConstraints
	var receiverMerkleProofs [operator.BatchSize]merkle.MerkleProof
	var transfersConstraints [operator.BatchSize]operator.TransferConstraints

	transfersRoot, err := r.TransfersRoot(transfers)
	if err != nil {
		return nil, fmt.Errorf("generate transactions: %w", err)
	}

	operatorReward := big.NewInt(0)
	for i := 0; i < len(transfers); i++ {
		transfersConstraints[i] = transfers[i].Constraints()

		fee := big.NewInt(0)
		transfers[i].Fee.BigInt(fee)
		reward := big.NewInt(0).Div(fee, big.NewInt(2))
		operatorReward.Add(operatorReward, reward)

		receiver, err := r.State.ReadAccount(uint64(i))
		if err != nil {
			return nil, fmt.Errorf("read account: %w", err)
		}

		receiverConstraints[i] = receiver.Constraints()

		root, receiverMerkleProof, err := r.State.MerkleProof(receiver.Index.Uint64())
		if err != nil {
			return nil, fmt.Errorf("create state: %w", err)
		}

		receiverMerkleProofs[i] = merkle.MerkleProof{
			RootHash: root,
			Path:     receiverMerkleProof,
		}

		err = r.UpdateReceiver(&receiver, &transfers[i])
		if err != nil {
			return nil, fmt.Errorf("failed to write account: %v", err)
		}
	}

	sourceOperator, err := r.State.ReadAccount(operator.NumberAccounts - 1)
	if err != nil {
		return nil, fmt.Errorf("read account: %w", err)
	}

	sourceOperatorConstraints := sourceOperator.Constraints()

	root, sourceOperatorMerkleProof, err := r.State.MerkleProof(sourceOperator.Index.Uint64())
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	sourceOperatorMerkleProofConstraints := merkle.MerkleProof{
		RootHash: root,
		Path:     sourceOperatorMerkleProof,
	}

	err = r.UpdateOperator(&sourceOperator, operatorReward)
	if err != nil {
		return nil, fmt.Errorf("failed to write account: %v", err)
	}

	targetOperator, err := r.State.ReadAccount(operator.NumberAccounts - 2)
	if err != nil {
		return nil, fmt.Errorf("read account: %w", err)
	}

	targetOperatorConstraints := targetOperator.Constraints()

	root, targetOperatorMerkleProof, err := r.State.MerkleProof(targetOperator.Index.Uint64())
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	targetOperatorMerkleProofConstraints := merkle.MerkleProof{
		RootHash: root,
		Path:     targetOperatorMerkleProof,
	}

	err = r.UpdateOperator(&targetOperator, operatorReward)
	if err != nil {
		return nil, fmt.Errorf("failed to write account: %v", err)
	}

	postStateRoot, err := r.State.Root()
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	assignment := &operator.ClaimCircuit{
		Receiver:                  receiverConstraints,
		SourceOperator:            sourceOperatorConstraints,
		TargetOperator:            targetOperatorConstraints,
		MerkleProofSourceOperator: sourceOperatorMerkleProofConstraints,
		MerkleProofTargetOperator: targetOperatorMerkleProofConstraints,
		MerkleProofReceiver:       receiverMerkleProofs,
		Transfers:                 transfersConstraints,
		PreStateRoot:              preStateRoot,
		PostStateRoot:             postStateRoot,
		TransactionsRoot:          transfersRoot,
	}

	w, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("failed to create witness: %v", err)
	}

	return w, nil
}

func (r *Rollup) TransfersRoot(transfers []operator.Transfer) ([]byte, error) {
	hFunc := mimc.NewMiMC()
	transferData := make([]byte, hFunc.Size()*len(transfers))

	for i := 0; i < len(transfers); i++ {
		copy(transferData[i*hFunc.Size():(i+1)*hFunc.Size()], transfers[i].Hash(hFunc))
	}

	var buf bytes.Buffer
	_, err := buf.Write(transferData)
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	root, err := merkletree.ReaderRoot(&buf, hFunc, hFunc.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer root: %v", err)
	}

	return root, nil
}

func (r *Rollup) UpdateSender(account *operator.Account, transfer *operator.Transfer) error {
	account.Nonce.Add(account.Nonce, big.NewInt(1))

	amount := big.NewInt(0)
	transfer.Amount.BigInt(amount)

	fee := big.NewInt(0)
	transfer.Fee.BigInt(fee)

	sum := big.NewInt(0).Add(amount, fee)
	account.Balance.Sub(account.Balance, sum)

	err := r.State.WriteAccount(*account)
	if err != nil {
		return fmt.Errorf("failed to write account: %v", err)
	}

	return nil
}

func (r *Rollup) UpdateReceiver(account *operator.Account, transfer *operator.Transfer) error {
	amount := big.NewInt(0)
	transfer.Amount.BigInt(amount)

	account.Balance.Add(account.Balance, amount)

	err := r.State.WriteAccount(*account)
	if err != nil {
		return fmt.Errorf("failed to write account: %v", err)
	}

	return nil
}

func (r *Rollup) UpdateOperator(account *operator.Account, reward *big.Int) error {
	account.Balance.Add(account.Balance, reward)

	err := r.State.WriteAccount(*account)
	if err != nil {
		return fmt.Errorf("failed to write account: %v", err)
	}

	return nil
}
