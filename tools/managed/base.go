package managed

import (
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/executable"
)

// BaseMCMSProposal is the base struct for all MCMS proposals
// Note: this type should never be utilized directly which is why it is private
type baseMCMSProposal struct {
	executable.ExecutableMCMSProposal

	// This is intended to be displayed as-is to signers, to give them
	// context for the change. File authors should templatize strings for
	// this purpose in their pipelines.
	Description string `json:"description"`

	// Operations to be executed
	Transactions []ChainOperation `json:"transactions"`
}

func (m *baseMCMSProposal) Validate() error {
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
		if _, ok := m.ChainMetadata[t.GetChainIdentifier()]; !ok {
			return &errors.ErrMissingChainDetails{
				ChainIdentifier: t.GetChainIdentifier(),
				Parameter:       "chain metadata",
			}
		}
	}
	return nil
}

func (m *baseMCMSProposal) AddSignature(sig executable.Signature) {
	m.Signatures = append(m.Signatures, sig)
}

func (m *baseMCMSProposal) ToExecutableMCMSProposal() executable.ExecutableMCMSProposal {
	raw := executable.ExecutableMCMSProposal{
		Version:              m.Version,
		ValidUntil:           m.ValidUntil,
		Signatures:           m.Signatures,
		OverridePreviousRoot: m.OverridePreviousRoot,
		Transactions:         make([]executable.ChainOperation, 0),
		ChainMetadata:        make(map[string]executable.ExecutableMCMSChainMetadata),
	}

	for k, v := range m.ChainMetadata {
		raw.ChainMetadata[k] = executable.ExecutableMCMSChainMetadata{
			NonceOffset: v.NonceOffset,
			MCMAddress:  v.MCMAddress,
		}
	}

	return raw
}
