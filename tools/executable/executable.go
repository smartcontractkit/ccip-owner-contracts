package executable

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

type ExecutableMCMSProposalBase struct {
	Version              string      `json:"version"`
	ValidUntil           string      `json:"validUntil"`
	Signatures           []Signature `json:"signatures"`
	OverridePreviousRoot bool        `json:"overridePreviousRoot"`

	// Map of chain identifier to chain metadata
	ChainMetadata map[string]ExecutableMCMSChainMetadata `json:"chainMetadata"`
}

type ExecutableMCMSChainMetadata struct {
	NonceOffset uint64         `json:"nonceOffset"`
	MCMAddress  common.Address `json:"mcmAddress"`
}

func (m ExecutableMCMSProposalBase) Validate() error {
	if m.Version == "" {
		return &errors.ErrInvalidVersion{
			ReceivedVersion: m.Version,
		}
	}

	if m.ValidUntil == "" {
		return &errors.ErrInvalidValidUntil{
			ReceivedValidUntil: m.ValidUntil,
		}
	}

	if len(m.ChainMetadata) == 0 {
		return &errors.ErrNoChainMetadata{}
	}

	return nil
}

type ExecutableMCMSProposal struct {
	ExecutableMCMSProposalBase

	// Operations to be executed
	Transactions []ChainOperation `json:"transactions"`
}

func (m *ExecutableMCMSProposal) Validate() error {
	if err := m.ExecutableMCMSProposalBase.Validate(); err != nil {
		return err
	}

	if len(m.Transactions) == 0 {
		return &errors.ErrNoTransactions{}
	}

	// Validate all chains in transactions have an entry in chain metadata
	for _, t := range m.Transactions {
		if _, ok := m.ChainMetadata[t.ChainIdentifier]; !ok {
			return &errors.ErrMissingChainDetails{
				ChainIdentifier: t.ChainIdentifier,
				Parameter:       "chain metadata",
			}
		}
	}

	return nil
}

func (m *ExecutableMCMSProposal) SigningHash(clients map[string]bind.ContractBackend) ([]byte, error) {
	tree, err := m.ConstructMerkleTree(clients)
	if err != nil {
		return nil, err
	}

	// convert validUntil to a big.Int
	validUntil, ok := new(big.Int).SetString(m.ValidUntil, 10)
	if !ok {
		return nil, &errors.ErrInvalidValidUntil{
			ReceivedValidUntil: m.ValidUntil,
		}
	}

	return append([]byte("\x19Ethereum Signed Message:\n"), crypto.Keccak256(tree.Root.Bytes(), validUntil.Bytes())...), nil
}

func (m *ExecutableMCMSProposal) ConstructMerkleTree(clients map[string]bind.ContractBackend) (*MerkleTree, error) {
	txCounts := calculateTransactionCounts(m.Transactions)
	currentOpCounts := m.GetCurrentOpCounts(clients)

	rootMetadatas, err := buildRootMetadatas(m.ChainMetadata, txCounts, currentOpCounts, m.OverridePreviousRoot)
	if err != nil {
		return nil, err
	}

	ops, err := buildOperations(m.Transactions, rootMetadatas, txCounts)
	if err != nil {
		return nil, err
	}

	chainIdentifiers := sortedChainIdentifiers(m.ChainMetadata)

	tree, err := buildMerkleTree(chainIdentifiers, rootMetadatas, ops)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func (m *ExecutableMCMSProposal) ValidateSignatures(clients map[string]bind.ContractBackend) error {
	hash, err := m.SigningHash(clients)
	if err != nil {
		return err
	}

	recoveredSigners := make([]common.Address, len(m.Signatures))
	for _, sig := range m.Signatures {
		recoveredAddr, err := recoverAddressFromSignature(hash, sig.ToBytes())
		if err != nil {
			return err
		}
		recoveredSigners = append(recoveredSigners, recoveredAddr)
	}

	wrappers, err := m.getAllMCMSWrappers(clients)
	if err != nil {
		return err
	}

	// Validate that all signers are valid
	for chain, wrapper := range wrappers {
		config, err := wrapper.GetConfig(&bind.CallOpts{})
		if err != nil {
			return err
		}

		for _, signer := range recoveredSigners {
			found := false
			for _, mcmsSigner := range config.Signers {
				if mcmsSigner.Addr == signer {
					found = true
					break
				}
			}

			if !found {
				return &errors.ErrInvalidSignature{
					ChainIdentifier:  chain,
					RecoveredAddress: signer,
				}
			}
		}
	}

	return nil
}

func (m *ExecutableMCMSProposal) GetCurrentOpCounts(clients map[string]bind.ContractBackend) map[string]big.Int {
	opCounts := make(map[string]big.Int)
	wrappers, err := m.getAllMCMSWrappers(clients)
	if err != nil {
		return nil
	}

	for chain, wrapper := range wrappers {
		opCount, err := wrapper.GetOpCount(&bind.CallOpts{})
		if err != nil {
			return nil
		}

		opCounts[chain] = *opCount
	}

	return opCounts
}

func (m *ExecutableMCMSProposal) getAllMCMSWrappers(clients map[string]bind.ContractBackend) (map[string]*gethwrappers.ManyChainMultiSig, error) {
	mcmsWrappers := make(map[string]*gethwrappers.ManyChainMultiSig)

	for chain, chainMetadata := range m.ChainMetadata {
		client, ok := clients[chain]
		if !ok {
			return nil, &errors.ErrMissingChainClient{
				ChainIdentifier: chain,
			}
		}

		mcms, err := gethwrappers.NewManyChainMultiSig(chainMetadata.MCMAddress, client)
		if err != nil {
			return nil, err
		}

		mcmsWrappers[chain] = mcms
	}

	return mcmsWrappers, nil
}
