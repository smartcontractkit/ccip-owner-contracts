package signing

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Define a sample struct
	type SampleStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	// Create a sample JSON file
	sampleData := SampleStruct{
		Name:  "John Doe",
		Age:   30,
		Email: "johndoe@example.com",
	}
	sampleJSON, err := json.Marshal(sampleData)
	assert.NoError(t, err)
	err = os.WriteFile(file.Name(), sampleJSON, 0644)
	assert.NoError(t, err)

	// Call the FromFile function
	var result SampleStruct
	err = FromFile(file.Name(), &result)
	assert.NoError(t, err)

	// Assert the result
	assert.Equal(t, sampleData, result)
}

func TestProposalFromFile(t *testing.T) {
	mcmsProposal := mcms.Proposal{
		Version:              "1",
		ValidUntil:           100,
		Signatures:           []mcms.Signature{},
		Transactions:         []mcms.ChainOperation{},
		OverridePreviousRoot: false,
		Description:          "Test Proposal",
		ChainMetadata:        make(map[mcms.ChainIdentifier]mcms.ChainMetadata),
	}

	tempFile, err := os.CreateTemp("", "mcms.json")
	assert.NoError(t, err)

	proposalBytes, err := json.Marshal(mcmsProposal)
	assert.NoError(t, err)
	err = os.WriteFile(tempFile.Name(), proposalBytes, 0644)
	assert.NoError(t, err)

	fileProposal, err := ProposalFromFile(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, mcmsProposal, *fileProposal)
}

func TestTimelockProposalFromFile(t *testing.T) {
	mcmsProposal := timelock.MCMSWithTimelockProposal{
		Proposal: mcms.Proposal{
			Version:              "1",
			ValidUntil:           100,
			Signatures:           []mcms.Signature{},
			OverridePreviousRoot: false,
			Description:          "Test Proposal",
		},
		ChainMetadata: make(map[mcms.ChainIdentifier]timelock.MCMSWithTimelockChainMetadata),
		Transactions:  make([]timelock.BatchChainOperation, 0),
		Operation:     timelock.Schedule,
		MinDelay:      "1h",
	}

	tempFile, err := os.CreateTemp("", "timelock.json")
	assert.NoError(t, err)

	proposalBytes, err := json.Marshal(mcmsProposal)
	assert.NoError(t, err)
	err = os.WriteFile(tempFile.Name(), proposalBytes, 0644)
	assert.NoError(t, err)

	fileProposal, err := TimelockProposalFromFile(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, mcmsProposal, *fileProposal)
}