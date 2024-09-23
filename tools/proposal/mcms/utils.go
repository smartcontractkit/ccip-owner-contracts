package mcms

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/configwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

type ContractDeployBackend interface {
	bind.ContractBackend
	bind.DeployBackend
}

func transformMCMAddresses(metadatas map[ChainIdentifier]ChainMetadata) map[ChainIdentifier]common.Address {
	m := make(map[ChainIdentifier]common.Address)
	for k, v := range metadatas {
		m[k] = v.MCMAddress
	}
	return m
}

func transformSignatures(signatures []Signature) []gethwrappers.ManyChainMultiSigSignature {
	sigs := make([]gethwrappers.ManyChainMultiSigSignature, len(signatures))
	for i, sig := range signatures {
		sigs[i] = sig.ToGethSignature()
	}
	return sigs
}

func transformHashes(hashes []common.Hash) [][32]byte {
	m := make([][32]byte, len(hashes))
	for i, h := range hashes {
		m[i] = [32]byte(h)
	}
	return m
}

func transformMCMSConfigs(configs map[ChainIdentifier]gethwrappers.ManyChainMultiSigConfig) (map[ChainIdentifier]*configwrappers.Config, error) {
	m := make(map[ChainIdentifier]*configwrappers.Config)
	for k, v := range configs {
		config, err := configwrappers.NewConfigFromRaw(v)
		if err != nil {
			return nil, err
		}
		m[k] = config
	}
	return m, nil
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

// Generic function to read a file and unmarshal its contents into the provided struct
func FromFile(filePath string, out interface{}) error {
	// Load file from path
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal JSON into the provided struct
	err = json.Unmarshal(fileBytes, out)
	if err != nil {
		return err
	}

	return nil
}

func WriteProposalToFile(proposal interface{}, filePath string) error {
	proposalBytes, err := json.Marshal(proposal)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, proposalBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Just run this locally to sign from the ledger.
func SignPlainKey(privateKey *ecdsa.PrivateKey, proposal MCMSProposal) error {
	// Validate proposal
	err := proposal.Validate()
	if err != nil {
		return err
	}

	executor, err := proposal.ToExecutor(false) // TODO: pass in a real backend
	if err != nil {
		return err
	}

	// Get the signing hash
	payload, err := executor.SigningHash()
	if err != nil {
		return err
	}

	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), privateKey)
	if err != nil {
		return err
	}

	// Unmarshal signature
	sigObj, err := NewSignatureFromBytes(sig)
	if err != nil {
		return err
	}

	// Add signature to proposal
	proposal.AddSignature(sigObj)
	return nil
}
