package signing

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/managed"
)

func ProposalFromFile(proposalType managed.MCMSProposalType, filePath string) (managed.MCMSProposal, error) {
	var out managed.MCMSProposal
	switch proposalType {
	case managed.MCMSOnly:
		out = &managed.MCMSOnlyProposal{}
	case managed.MCMSWithTimelock:
		out = &managed.MCMSWithTimelockProposal{}
	case managed.MCMSWithMultisend:
		out = &managed.MCMSWithMultisendProposal{}
	default:
		return nil, &errors.ErrInvalidProposalType{ReceivedProposalType: string(proposalType)}
	}

	json.Unmarshal([]byte(filePath), out)
	return out, nil
}

func WriteProposalToFile(proposal managed.MCMSProposal, filePath string) error {
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
