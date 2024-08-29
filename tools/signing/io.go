package signing

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/mcms_proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/timelock_proposal"
)

func ProposalFromFile(filePath string) (*mcms_proposal.Proposal, error) {
	var out mcms_proposal.Proposal

	// Load file from path
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fileBytes, out)
	return &out, nil
}

func TimelockProposalFromFile(filePath string) (*timelock_proposal.MCMSWithTimelockProposal, error) {
	var out timelock_proposal.MCMSWithTimelockProposal

	// Load file from path
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fileBytes, out)
	return &out, nil
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
