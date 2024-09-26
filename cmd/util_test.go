package cmd

import (
	"crypto/ecdsa"
	"encoding/json"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/stretchr/testify/assert"
)

func TestLoadPrivateKey(t *testing.T) {
	// Set up test environment
	privateKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	// Write private key to .env file
	err := os.WriteFile(".env", []byte("PRIVATE_KEY="+privateKey), 0644)
	assert.NoError(t, err)

	// Call the function
	pk, err := LoadPrivateKey()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, pk)
	assert.IsType(t, &ecdsa.PrivateKey{}, pk)
}

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
	mcmsProposal := mcms.MCMSProposal{
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
		MCMSProposal: mcms.MCMSProposal{
			Version:              "1",
			ValidUntil:           100,
			Signatures:           []mcms.Signature{},
			OverridePreviousRoot: false,
			Description:          "Test Proposal",
			ChainMetadata:        make(map[mcms.ChainIdentifier]mcms.ChainMetadata),
		},
		TimelockAddresses: make(map[mcms.ChainIdentifier]common.Address),
		Transactions:      make([]timelock.BatchChainOperation, 0),
		Operation:         timelock.Schedule,
		MinDelay:          "1h",
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
