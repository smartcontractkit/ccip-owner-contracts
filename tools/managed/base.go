package managed

import (
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/executable"
)

type BaseMCMSProposal struct {
	executable.ExecutableMCMSProposalBase

	// This is intended to be displayed as-is to signers, to give them
	// context for the change. File authors should templatize strings for
	// this purpose in their pipelines.
	Description string `json:"description"`

	// Map of chain identifier to chain metadata
	ChainMetadata map[string]executable.ExecutableMCMSChainMetadata `json:"chainMetadata"`

	// Operations to be executed
	Transactions []ChainOperation `json:"transactions"`
}

func (m BaseMCMSProposal) Validate() error {
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

	if m.Description == "" {
		return &errors.ErrInvalidDescription{
			ReceivedDescription: m.Description,
		}
	}

	if len(m.ChainMetadata) == 0 {
		return &errors.ErrNoChainMetadata{}
	}

	if len(m.Transactions) == 0 {
		return &errors.ErrNoTransactions{}
	}

	// Validate all chains in transactions have an entry in chain metadata
	for _, t := range m.Transactions {
		if _, ok := m.ChainMetadata[t.GetChainIdentifier()]; !ok {
			return &errors.ErrMissingChainMetadata{
				ChainIdentifier: t.GetChainIdentifier(),
			}
		}
	}
	return nil
}

func (m BaseMCMSProposal) AddSignature(sig executable.Signature) error {
	m.Signatures = append(m.Signatures, sig)
	return nil
}

func (m BaseMCMSProposal) ToExecutableMCMSProposal() executable.ExecutableMCMSProposal {
	raw := executable.ExecutableMCMSProposal{
		ExecutableMCMSProposalBase: m.ExecutableMCMSProposalBase,
		ChainMetadata:              make(map[string]executable.ExecutableMCMSChainMetadata),
		Transactions:               make([]executable.ChainOperation, 0),
	}

	for k, v := range m.ChainMetadata {
		raw.ChainMetadata[k] = executable.ExecutableMCMSChainMetadata{
			NonceOffset: v.NonceOffset,
			MCMAddress:  v.MCMAddress,
		}
	}

	return raw
}
