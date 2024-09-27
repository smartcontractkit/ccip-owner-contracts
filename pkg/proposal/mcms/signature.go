package mcms

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
)

type Signature struct {
	R common.Hash
	S common.Hash
	V uint8
}

func NewSignatureFromBytes(sig []byte) (Signature, error) {
	if len(sig) != 65 {
		return Signature{}, fmt.Errorf("invalid signature length: %d", len(sig))
	}

	return Signature{
		R: common.BytesToHash(sig[:32]),
		S: common.BytesToHash(sig[32:64]),
		V: uint8(sig[64]),
	}, nil
}

func (s Signature) ToGethSignature() gethwrappers.ManyChainMultiSigSignature {
	if s.V < 2 {
		s.V += 27
	}

	return gethwrappers.ManyChainMultiSigSignature{
		R: [32]byte(s.R.Bytes()),
		S: [32]byte(s.S.Bytes()),
		V: s.V,
	}
}

func (s Signature) ToBytes() []byte {
	return append(s.R.Bytes(), append(s.S.Bytes(), []byte{byte(s.V)}...)...)
}

func (s Signature) Recover(hash common.Hash) (common.Address, error) {
	return recoverAddressFromSignature(hash, s.ToBytes())
}

func recoverAddressFromSignature(hash common.Hash, sig []byte) (common.Address, error) {
	// The signature should be 65 bytes, and the last byte is the recovery id (v).
	if len(sig) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length")
	}

	// Adjust the recovery id (v) if needed. Ethereum signatures expect 27 or 28.
	// But `crypto.SigToPub` expects 0 or 1.
	if sig[64] > 1 {
		sig[64] -= 27
	}

	// Recover the public key from the signature and the message hash
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return common.Address{}, err
	}

	// Derive the Ethereum address from the public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, nil
}
