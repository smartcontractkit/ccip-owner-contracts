package signing

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
)

func ProposalFromFile(filePath string) (*mcms.Proposal, error) {
	var out mcms.Proposal

	// Load file from path
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fileBytes, out)
	return &out, nil
}

func TimelockProposalFromFile(filePath string) (*timelock.MCMSWithTimelockProposal, error) {
	var out timelock.MCMSWithTimelockProposal

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
