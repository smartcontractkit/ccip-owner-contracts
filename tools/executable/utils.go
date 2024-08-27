package executable

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/configwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

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
		m[i] = [32]byte(h)
	}
	return m
}

func mapMCMSConfigs(configs map[string]gethwrappers.ManyChainMultiSigConfig) map[string]configwrappers.Config {
	m := make(map[string]configwrappers.Config)
	for k, v := range configs {
		m[k] = *configwrappers.NewConfigFromRaw(v)
	}
	return m
}
