package mcms_proposal

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var TestAddress = common.HexToAddress("0x1234567890abcdef")
var TestChain = "chain1"

func TestMCMSOnlyProposal_Validate_Success(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain,
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.NoError(t, err)
}

func TestMCMSOnlyProposal_Validate_InvalidVersion(t *testing.T) {
	proposal := &Proposal{
		Version:              "",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain,
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid version: ")
}

func TestMCMSOnlyProposal_AddSignature(t *testing.T) {
	proposal := Proposal{
		Version:    "1.0.0",
		ValidUntil: 2004259681,
		Signatures: []Signature{},
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			ChainOperation{
				ChainIdentifier: TestChain,
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	sig := Signature{
		R: common.HexToHash("0x1234567890abcdef"),
		S: common.HexToHash("0x1234567890abcdef"),
		V: 27,
	}

	proposal.AddSignature(sig)
	assert.Len(t, proposal.Signatures, 1)
	assert.Equal(t, proposal.Signatures[0], sig)
}

func TestMCMSOnlyProposal_Validate_InvalidValidUntil(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           0,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain,
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid valid until: 0")
}

func TestMCMSOnlyProposal_Validate_InvalidChainMetadata(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata:        map[string]ChainMetadata{},
		Description:          "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain,
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no chain metadata")
}

func TestMCMSOnlyProposal_Validate_InvalidDescription(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain,
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid description: ")
}

func TestMCMSOnlyProposal_Validate_NoTransactions(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		Description:          "Sample description",
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Transactions: []ChainOperation{},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no transactions")
}

func TestMCMSOnlyProposal_Validate_MissingChainMetadataForTransaction(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[string]ChainMetadata{
			TestChain: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: "chain2",
				Operation: Operation{
					To:           TestAddress,
					Value:        0,
					Data:         "0x",
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "missing chain metadata for chain chain2")
}
