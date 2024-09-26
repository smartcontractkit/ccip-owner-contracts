package mcms

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/errors"
)

type ChainMetadata struct {
	StartingOpCount uint64         `json:"startingOpCount"`
	MCMAddress      common.Address `json:"mcmAddress"`
}

// MCMSProposal is a struct where the target contract is an MCMS contract
// with no forwarder contracts. This type does not support any type of atomic contract
// call batching, as the MCMS contract natively doesn't support batching
type MCMSProposal struct {
	Version              string      `json:"version"`
	ValidUntil           uint32      `json:"validUntil"`
	Signatures           []Signature `json:"signatures"`
	OverridePreviousRoot bool        `json:"overridePreviousRoot"`

	// Map of chain identifier to chain metadata
	ChainMetadata map[ChainIdentifier]ChainMetadata `json:"chainMetadata"`

	// This is intended to be displayed as-is to signers, to give them
	// context for the change. File authors should templatize strings for
	// this purpose in their pipelines.
	Description string `json:"description"`

	// Operations to be executed
	Transactions []ChainOperation `json:"transactions"`
}

func NewProposal(
	version string,
	validUntil uint32,
	signatures []Signature,
	overridePreviousRoot bool,
	chainMetadata map[ChainIdentifier]ChainMetadata,
	description string,
	transactions []ChainOperation,
) (*MCMSProposal, error) {
	proposal := MCMSProposal{
		Version:              version,
		ValidUntil:           validUntil,
		Signatures:           signatures,
		OverridePreviousRoot: overridePreviousRoot,
		ChainMetadata:        chainMetadata,
		Description:          description,
		Transactions:         transactions,
	}

	err := proposal.Validate()
	if err != nil {
		return nil, err
	}

	return &proposal, nil
}

func (m *MCMSProposal) Validate() error {
	if m.Version == "" {
		return &errors.ErrInvalidVersion{
			ReceivedVersion: m.Version,
		}
	}

	// Get the current Unix timestamp as an int64
	currentTime := time.Now().Unix()

	if m.ValidUntil <= uint32(currentTime) {
		// ValidUntil is a Unix timestamp, so it should be greater than the current time
		return &errors.ErrInvalidValidUntil{
			ReceivedValidUntil: m.ValidUntil,
		}
	}

	if len(m.ChainMetadata) == 0 {
		return &errors.ErrNoChainMetadata{}
	}

	if len(m.Transactions) == 0 {
		return &errors.ErrNoTransactions{}
	}

	if m.Description == "" {
		return &errors.ErrInvalidDescription{
			ReceivedDescription: m.Description,
		}
	}

	// Validate all chains in transactions have an entry in chain metadata
	for _, t := range m.Transactions {
		if _, ok := m.ChainMetadata[t.ChainIdentifier]; !ok {
			return &errors.ErrMissingChainDetails{
				ChainIdentifier: uint64(t.ChainIdentifier),
				Parameter:       "chain metadata",
			}
		}
	}
	return nil
}

func (m *MCMSProposal) ToExecutor(sim bool) (*Executor, error) {
	// Create a new executor
	executor, err := NewProposalExecutor(m, sim)
	if err != nil {
		return nil, err
	}

	return executor, nil
}

func (m *MCMSProposal) AddSignature(signature Signature) {
	m.Signatures = append(m.Signatures, signature)
}
