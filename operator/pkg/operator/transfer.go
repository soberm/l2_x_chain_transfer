package operator

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	eddsa2 "github.com/consensys/gnark/std/signature/eddsa"
	"hash"
)

type Transfer struct {
	Nonce          uint64
	Amount         fr.Element
	SenderPubKey   eddsa.PublicKey
	ReceiverPubKey eddsa.PublicKey
	Fee            fr.Element
	Signature      eddsa.Signature
	Destination    fr.Element
}

func NewTransfer(amount, fee uint64, from, to eddsa.PublicKey, nonce uint64, destination uint64) Transfer {

	var res Transfer

	res.Nonce = nonce
	res.Amount.SetUint64(amount)
	res.Fee.SetUint64(fee)
	res.SenderPubKey = from
	res.ReceiverPubKey = to
	res.Destination.SetUint64(destination)

	return res
}

func (t *Transfer) Constraints() TransferConstraints {
	var senderPubKey eddsa2.PublicKey
	var receiverPubKey eddsa2.PublicKey
	var sig eddsa2.Signature

	senderPubKey.Assign(tedwards.BN254, t.SenderPubKey.Bytes())
	receiverPubKey.Assign(tedwards.BN254, t.ReceiverPubKey.Bytes())
	sig.Assign(tedwards.BN254, t.Signature.Bytes())

	return TransferConstraints{
		Nonce:          t.Nonce,
		Fee:            t.Fee,
		Amount:         t.Amount,
		SenderPubKey:   senderPubKey,
		ReceiverPubKey: receiverPubKey,
		Signature:      sig,
		Destination:    t.Destination,
	}
}

func (t *Transfer) Hash(h hash.Hash) []byte {
	h.Reset()
	var frNonce fr.Element

	frNonce.SetUint64(t.Nonce)
	b := frNonce.Bytes()
	_, _ = h.Write(b[:])
	b = t.Amount.Bytes()
	_, _ = h.Write(b[:])
	b = t.Fee.Bytes()
	_, _ = h.Write(b[:])
	b = t.SenderPubKey.A.X.Bytes()
	_, _ = h.Write(b[:])
	b = t.SenderPubKey.A.Y.Bytes()
	_, _ = h.Write(b[:])
	b = t.ReceiverPubKey.A.X.Bytes()
	_, _ = h.Write(b[:])
	b = t.ReceiverPubKey.A.Y.Bytes()
	_, _ = h.Write(b[:])
	b = t.Destination.Bytes()
	_, _ = h.Write(b[:])

	return h.Sum([]byte{})
}

func (t *Transfer) Sign(priv eddsa.PrivateKey, h hash.Hash) (eddsa.Signature, []byte, error) {

	h.Reset()
	msg := t.Hash(h)

	sigBin, err := priv.Sign(msg, h)
	if err != nil {
		return eddsa.Signature{}, nil, err
	}
	var sig eddsa.Signature
	if _, err := sig.SetBytes(sigBin); err != nil {
		return eddsa.Signature{}, nil, err
	}
	t.Signature = sig
	return sig, msg, nil
}
