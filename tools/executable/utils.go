package executable

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

func mapMCMSConfigs(configs map[string]gethwrappers.ManyChainMultiSigConfig) map[string]*configwrappers.Config {
	m := make(map[string]*configwrappers.Config)
	for k, v := range configs {
		m[k] = configwrappers.NewConfigFromRaw(v)
	}
	return m
}

// ABIEncode is the equivalent of abi.encode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIEncode(abiStr string, values ...interface{}) ([]byte, error) {
	// Create a dummy method with arguments
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "inputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}
	return res[4:], nil
}

// ABIDecode is the equivalent of abi.decode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIDecode(abiStr string, data []byte) ([]interface{}, error) {
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "outputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	return inAbi.Unpack("method", data)
}
