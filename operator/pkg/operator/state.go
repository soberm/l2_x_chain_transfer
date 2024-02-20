package operator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"hash"
	"sync"
)

type State struct {
	sync.RWMutex
	hFunc hash.Hash
	data  []byte
	hData []byte
}

func NewState(hFunc hash.Hash, accounts []*Account) (*State, error) {

	data := make([]byte, accountSize*len(accounts))
	hData := make([]byte, hFunc.Size()*len(accounts))

	for i, account := range accounts {

		hFunc.Reset()

		accountData := account.Serialize()
		_, err := hFunc.Write(accountData)
		if err != nil {
			return nil, fmt.Errorf("hash account: %w", err)
		}
		hashedAccountData := hFunc.Sum(nil)

		copy(data[i*accountSize:(i+1)*accountSize], accountData)
		copy(hData[i*hFunc.Size():(i+1)*hFunc.Size()], hashedAccountData)
	}

	return &State{
		hFunc: hFunc,
		data:  data,
		hData: hData,
	}, nil
}

func (s *State) ReadAccount(i uint64) (Account, error) {
	s.RLock()
	defer s.RUnlock()

	var res Account
	res.Deserialize(s.data[int(i)*accountSize : int(i)*accountSize+accountSize])
	return res, nil
}

func (s *State) WriteAccount(account Account) error {
	s.Lock()
	defer s.Unlock()

	i := int(account.Index.Int64())
	accountData := account.Serialize()

	copy(s.data[i*accountSize:], accountData)

	s.hFunc.Reset()
	_, err := s.hFunc.Write(accountData)
	if err != nil {
		return fmt.Errorf("hash account: %w", err)
	}
	copy(s.hData[i*s.hFunc.Size():(i+1)*s.hFunc.Size()], s.hFunc.Sum(nil))

	return nil
}

func (s *State) Root() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	var buf bytes.Buffer
	_, err := buf.Write(s.hData)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return merkletree.ReaderRoot(&buf, s.hFunc, s.hFunc.Size())
}

func (s *State) MerkleProof(i uint64) ([]byte, [][]byte, error) {
	s.RLock()
	defer s.RUnlock()

	var stateBuf bytes.Buffer
	_, err := stateBuf.Write(s.hData)
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	root, proof, numLeaves, _ := merkletree.BuildReaderProof(&stateBuf, s.hFunc, s.hFunc.Size(), i)

	if !merkletree.VerifyProof(s.hFunc, root, proof, i, numLeaves) {
		return nil, proof, errors.New("invalid merkle proof")
	}

	return root, proof, nil
}
