package executable

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

var MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP = crypto.Keccak256Hash([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP"))
var MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA = crypto.Keccak256Hash([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA"))

type ContractDeployBackend interface {
	bind.ContractBackend
	bind.DeployBackend
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

func mapMCMAddresses(metadatas map[string]ExecutableMCMSChainMetadata) map[string]common.Address {
	m := make(map[string]common.Address)
	for k, v := range metadatas {
		m[k] = v.MCMAddress
	}
	return m
}

func mapSignatures(signatures []Signature) []gethwrappers.ManyChainMultiSigSignature {
	sigs := make([]gethwrappers.ManyChainMultiSigSignature, len(signatures))
	for i, sig := range signatures {
		sigs[i] = sig.ToGethSignature()
	}
	return sigs
}

func mapHashes(hashes []common.Hash) [][32]byte {
	m := make([][32]byte, len(hashes))
	for i, h := range hashes {
		m[i] = hashToBytes32(h)
	}
	return m
}

func hashToBytes32(hash common.Hash) [32]byte {
	return [32]byte(hash)
}
