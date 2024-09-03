package signing

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/mcmsproposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/timelockproposal"
)

func ProposalFromFile(filePath string) (*mcmsproposal.Proposal, error) {
	var out mcmsproposal.Proposal

	// Load file from path
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fileBytes, out)
	return &out, nil
}

func TimelockProposalFromFile(filePath string) (*timelockproposal.MCMSWithTimelockProposal, error) {
	var out timelockproposal.MCMSWithTimelockProposal

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
