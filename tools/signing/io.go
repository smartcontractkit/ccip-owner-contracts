package signing

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
)

// Generic function to read a file and unmarshal its contents into the provided struct
func FromFile(filePath string, out interface{}) error {
	// Load file from path
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal JSON into the provided struct
	err = json.Unmarshal(fileBytes, out)
	if err != nil {
		return err
	}

	return nil
}

// Usage for mcms.Proposal
func ProposalFromFile(filePath string) (*mcms.Proposal, error) {
	var out mcms.Proposal
	err := FromFile(filePath, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Usage for timelock.MCMSWithTimelockProposal
func TimelockProposalFromFile(filePath string) (*timelock.MCMSWithTimelockProposal, error) {
	var out timelock.MCMSWithTimelockProposal
	err := FromFile(filePath, &out)
	if err != nil {
		return nil, err
	}
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
