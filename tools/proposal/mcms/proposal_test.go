package mcms

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
	proposal, err := NewProposal(
		"1.0",
		2004259681,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		"Sample description",
		[]ChainOperation{
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
	)

	assert.NoError(t, err)
	assert.NotNil(t, proposal)
}

func TestMCMSOnlyProposal_Validate_InvalidVersion(t *testing.T) {
	proposal, err := NewProposal(
		"",
		2004259681,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		"Sample description",
		[]ChainOperation{
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
	)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid version: ")
	assert.Nil(t, proposal)
}

func TestMCMSOnlyProposal_Validate_InvalidValidUntil(t *testing.T) {
	proposal, err := NewProposal(
		"1.0",
		0,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		"Sample description",
		[]ChainOperation{
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
	)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid valid until: 0")
	assert.Nil(t, proposal)
}

func TestMCMSOnlyProposal_Validate_InvalidChainMetadata(t *testing.T) {
	proposal, err := NewProposal(
		"1.0",
		2004259681,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{},
		"Sample description",
		[]ChainOperation{
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
	)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no chain metadata")
	assert.Nil(t, proposal)
}

func TestMCMSOnlyProposal_Validate_InvalidDescription(t *testing.T) {
	proposal, err := NewProposal(
		"1.0",
		2004259681,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		"",
		[]ChainOperation{
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
	)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid description: ")
	assert.Nil(t, proposal)
}

func TestMCMSOnlyProposal_Validate_NoTransactions(t *testing.T) {
	proposal, err := NewProposal(
		"1.0",
		2004259681,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		"Sample description",
		[]ChainOperation{},
	)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no transactions")
	assert.Nil(t, proposal)
}

func TestMCMSOnlyProposal_Validate_MissingChainMetadataForTransaction(t *testing.T) {
	proposal, err := NewProposal(
		"1.0",
		2004259681,
		[]Signature{},
		false,
		map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		"Sample description",
		[]ChainOperation{
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
	)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "missing chain metadata for chain 3")
	assert.Nil(t, proposal)
}
