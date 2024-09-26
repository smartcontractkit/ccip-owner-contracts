package cmd

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
)

func LoadPrivateKey() (*ecdsa.PrivateKey, error) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	// Load PrivateKey
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		return nil, errors.New("PRIVATE_KEY not found in .env file")
	}

	// Convert to ecdsa
	ecdsa, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}
	return ecdsa, nil
}

func LoadProposal(proposalType proposal.ProposalType, filePath string) (proposal.Proposal, error) {
	switch proposalType {
	case proposal.MCMS:
		return ProposalFromFile(filePath)
	case proposal.MCMSWithTimelock:
		return TimelockProposalFromFile(filePath)
	default:
		return nil, errors.New("unknown proposal type")
	}
}

// Usage for mcms.MCMSProposal
func ProposalFromFile(filePath string) (*mcms.MCMSProposal, error) {
	var out mcms.MCMSProposal
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
