package managed

import (
	"github.com/smartcontractkit/ccip-owner-contracts/tools/executable"
)

type MCMSOnlyChainMetadata struct {
	executable.ExecutableMCMSChainMetadata
}

type MCMSOnlyProposal struct {
	BaseMCMSProposal

	// Operations to be executed
	Transactions []DetailedChainOperation `json:"transactions"`
}

func (m MCMSOnlyProposal) ToExecutableMCMSProposal() (executable.ExecutableMCMSProposal, error) {
	raw := m.BaseMCMSProposal.ToExecutableMCMSProposal()

	for _, t := range m.Transactions {
		raw.Transactions = append(raw.Transactions, executable.ChainOperation{
			ChainIdentifier: t.ChainIdentifier,
			Operation:       t.Operation,
		})
	}

	return raw, nil
}
