package managed

import (
	"github.com/smartcontractkit/ccip-owner-contracts/tools/executable"
)

type MCMSWithMultisendChainMetadata struct {
	executable.ExecutableMCMSChainMetadata
	MultisendAddress string `json:"multisendAddress"`
}

type MCMSWithMultisendProposal struct {
	baseMCMSProposal

	// Map of chain identifier to chain metadata
	ChainMetadata map[string]MCMSWithMultisendChainMetadata `json:"chainMetadata"`

	// Operations to be executed
	Transactions []DetailedBatchChainOperation `json:"transactions"`
}

func (m *MCMSWithMultisendProposal) ToExecutableMCMSProposal() (executable.ExecutableMCMSProposal, error) {
	raw := m.baseMCMSProposal.ToExecutableMCMSProposal()

	for _, t := range m.Transactions {
		raw.Transactions = append(raw.Transactions, executable.ChainOperation{
			ChainIdentifier: t.ChainIdentifier,

			// TODO: wrap the batch (or the batch) in a multisend:<operation> call here
			Operation: *&t.Batch[0].Operation,
		})
	}

	return raw, nil
}
