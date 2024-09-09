package signing

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
)

func LoadProposal(proposalType proposal.ProposalType, filePath string) (*mcms.Proposal, error) {
	switch proposalType {
	case proposal.MCMS:
		return ProposalFromFile(filePath)
	case proposal.MCMSWithTimelock:
		proposal, err := TimelockProposalFromFile(filePath)
		if err != nil {
			return nil, err
		}

		mcmsOnlyProposal, err := proposal.ToMCMSOnlyProposal()
		if err != nil {
			return nil, err
		}

		return &mcmsOnlyProposal, nil
	default:
		return nil, errors.New("unknown proposal type")
	}
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
