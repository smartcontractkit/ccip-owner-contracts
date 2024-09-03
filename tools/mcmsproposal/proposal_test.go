package mcmsproposal

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var TestAddress = common.HexToAddress("0x1234567890abcdef")
var TestChain1 = ChainIdentifier(3379446385462418246)
var TestChain2 = ChainIdentifier(16015286601757825753)
var TestChain3 = ChainIdentifier(10344971235874465080)

func TestMCMSOnlyProposal_Validate_Success(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:           TestAddress,
					Value:        big.NewInt(0),
					Data:         common.Hex2Bytes("0x"),
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
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:           TestAddress,
					Value:        big.NewInt(0),
					Data:         common.Hex2Bytes("0x"),
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

func TestMCMSOnlyProposal_Validate_InvalidValidUntil(t *testing.T) {
	proposal := &Proposal{
		Version:              "1.0",
		ValidUntil:           0,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:           TestAddress,
					Value:        big.NewInt(0),
					Data:         common.Hex2Bytes("0x"),
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
		ChainMetadata:        map[ChainIdentifier]ChainMetadata{},
		Description:          "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:           TestAddress,
					Value:        big.NewInt(0),
					Data:         common.Hex2Bytes("0x"),
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
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:           TestAddress,
					Value:        big.NewInt(0),
					Data:         common.Hex2Bytes("0x"),
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
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
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
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 1,
				MCMAddress:  TestAddress,
			},
		},
		Description: "Sample description",
		Transactions: []ChainOperation{
			{
				ChainIdentifier: 3,
				Operation: Operation{
					To:           TestAddress,
					Value:        big.NewInt(0),
					Data:         common.Hex2Bytes("0x"),
					ContractType: "Sample contract",
					Tags:         []string{"tag1", "tag2"},
				},
			},
		},
	}

	err := proposal.Validate()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "missing chain metadata for chain 3")
}
