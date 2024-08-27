package executable

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

type ContractDeployBackend interface {
	bind.ContractBackend
	bind.DeployBackend
}

type Operation struct {
	To    common.Address `json:"to"`
	Data  string         `json:"data"`
	Value uint64         `json:"value"`
}

type ChainOperation struct {
	ChainIdentifier string `json:"chainSelector"`
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
		m[i] = [32]byte(h)
	}
	return m
}
