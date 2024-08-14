package executable

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP = crypto.Keccak256([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP"))
var MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA = crypto.Keccak256([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA"))

type Signature struct {
	R string
	S string
	V uint8
}

func NewSignatureFromBytes(sig []byte) Signature {
	return Signature{
		R: common.Bytes2Hex(sig[:32]),
		S: common.Bytes2Hex(sig[32:64]),
		V: uint8(sig[64]),
	}
}

func (s Signature) ToBytes() []byte {
	return append(common.FromHex(s.R), append(common.FromHex(s.S), []byte{byte(s.V)}...)...)
}

func (s Signature) Recover(hash []byte) (common.Address, error) {
	return recoverAddressFromSignature(hash, s.ToBytes())
}

type Operation struct {
	To    common.Address
	Data  string
	Value uint64
}

type ChainOperation struct {
	ChainIdentifier string
	Operation
}

func recoverAddressFromSignature(hash, sig []byte) (common.Address, error) {
	// The signature should be 65 bytes, and the last byte is the recovery id (v).
	if len(sig) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length")
	}

	// Adjust the recovery id (v) if needed. Ethereum signatures expect 27 or 28.
	// But `crypto.SigToPub` expects 0 or 1.
	sig[64] -= 27

	// Recover the public key from the signature and the message hash
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return common.Address{}, err
	}

	// Derive the Ethereum address from the public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr, nil
}
