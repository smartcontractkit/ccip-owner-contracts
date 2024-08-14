package executable

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
			return &errors.ErrMissingChainMetadata{
				ChainIdentifier: t.ChainIdentifier,
			}
		}
	}

	return nil
}

func (m *ExecutableMCMSProposal) SigningHash() ([]byte, error) {
	_, root, err := m.ConstructMerkleTree()
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

	return append([]byte("\x19Ethereum Signed Message:\n"), crypto.Keccak256(common.FromHex(root), validUntil.Bytes())...), nil
}

func (m *ExecutableMCMSProposal) ConstructMerkleTree() (map[string]string, string, error) {
	txCounts := calculateTransactionCounts(m.Transactions)

	rootMetadatas, err := buildRootMetadatas(m.ChainMetadata, txCounts, m.OverridePreviousRoot)
	if err != nil {
		return nil, "", err
	}

	ops, err := buildOperations(m.Transactions, rootMetadatas, txCounts)
	if err != nil {
		return nil, "", err
	}

	chainIdentifiers := sortedChainIdentifiers(m.ChainMetadata)

	tree, root, err := buildMerkleTree(chainIdentifiers, ops)
	if err != nil {
		return nil, "", err
	}

	return tree, root, nil
}

func (m *ExecutableMCMSProposal) ValidateSignatures(clients map[string]ethclient.Client) error {
	hash, err := m.SigningHash()
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

func (m *ExecutableMCMSProposal) GetCurrentOpCounts(clients map[string]ethclient.Client) map[string]big.Int {
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

func (m *ExecutableMCMSProposal) getAllMCMSWrappers(clients map[string]ethclient.Client) (map[string]*gethwrappers.ManyChainMultiSig, error) {
	mcmsWrappers := make(map[string]*gethwrappers.ManyChainMultiSig)

	for chain, chainMetadata := range m.ChainMetadata {
		client, ok := clients[chain]
		if !ok {
			return nil, &errors.ErrMissingChainClient{
				ChainIdentifier: chain,
			}
		}

		mcms, err := gethwrappers.NewManyChainMultiSig(chainMetadata.MCMAddress, &client)
		if err != nil {
			return nil, err
		}

		mcmsWrappers[chain] = mcms
	}

	return mcmsWrappers, nil
}
