package configwrappers

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	config, err := NewConfig(1, signers, groupSigners)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, uint8(1), config.Quorum)
	assert.Equal(t, signers, config.Signers)
	assert.Equal(t, groupSigners, config.GroupSigners)
}

func TestNewConfigFromRaw(t *testing.T) {
	rawConfig := gethwrappers.ManyChainMultiSigConfig{
		GroupQuorums: [32]uint8{1, 1},
		GroupParents: [32]uint8{0, 0},
		Signers: []gethwrappers.ManyChainMultiSigSigner{
			{Addr: common.HexToAddress("0x1"), Group: 0},
			{Addr: common.HexToAddress("0x2"), Group: 1},
		},
	}
	config, err := NewConfigFromRaw(rawConfig)

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, uint8(1), config.Quorum)
	assert.Equal(t, []common.Address{common.HexToAddress("0x1")}, config.Signers)
	assert.Equal(t, uint8(1), config.GroupSigners[0].Quorum)
	assert.Equal(t, []common.Address{common.HexToAddress("0x2")}, config.GroupSigners[0].Signers)
}

func TestValidate_Success(t *testing.T) {
	// Test case 1: Valid configuration
	config, err := NewConfig(2, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}, []Config{})
	assert.NotNil(t, config)
	assert.NoError(t, err)
}

func TestValidate_InvalidQuorum(t *testing.T) {
	// Test case 2: Quorum is 0
	config, err := NewConfig(0, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}, []Config{})
	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Equal(t, "invalid MCMS config: Quorum must be greater than 0", err.Error())
}

func TestValidate_InvalidSigners(t *testing.T) {
	// Test case 3: No signers or groups
	config, err := NewConfig(2, []common.Address{}, []Config{})
	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Equal(t, "invalid MCMS config: Config must have at least one signer or group", err.Error())
}

func TestValidate_InvalidQuorumCount(t *testing.T) {
	// Test case 4: Quorum is greater than the number of signers and groups
	config, err := NewConfig(3, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}, []Config{})
	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Equal(t, "invalid MCMS config: Quorum must be less than or equal to the number of signers and groups", err.Error())
}

func TestValidate_InvalidGroupSigner(t *testing.T) {
	// Test case 5: Invalid group signer
	config, err := NewConfig(2, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}, []Config{
		{Quorum: 0, Signers: []common.Address{}},
	})
	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Equal(t, "invalid MCMS config: Quorum must be greater than 0", err.Error())
}

func TestToRawConfig(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	config, err := NewConfig(1, signers, groupSigners)
	assert.NotNil(t, config)
	assert.NoError(t, err)

	rawConfig := config.ToRawConfig()

	assert.Equal(t, [32]uint8{1, 1}, rawConfig.GroupQuorums)
	assert.Equal(t, [32]uint8{0, 0}, rawConfig.GroupParents)
	assert.Equal(t, common.HexToAddress("0x1"), rawConfig.Signers[0].Addr)
	assert.Equal(t, common.HexToAddress("0x2"), rawConfig.Signers[1].Addr)
	assert.Equal(t, common.HexToAddress("0x3"), rawConfig.Signers[2].Addr)
}

// Test case 0: Valid configuration with no signers or groups
// Configuration:
// Quorum: 0
// Signers: []
// Group signers: []
func TestExtractSetConfigInputs_EmptyConfig(t *testing.T) {
	config, err := NewConfig(0, []common.Address{}, []Config{})
	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Equal(t, "invalid MCMS config: Quorum must be greater than 0", err.Error())
}

// Test case 1: Valid configuration with some root signers and some groups
// Configuration:
// Quorum: 2
// Signers: [0x1, 0x2]
//
//	Group signers: [{
//		Quorum: 1
//		Signers: [0x3]
//		Group signers: []
//	}]
func TestExtractSetConfigInputs(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	config, err := NewConfig(2, signers, groupSigners)
	assert.NotNil(t, config)
	assert.NoError(t, err)

	groupQuorums, groupParents, signerAddresses, signerGroups := config.ExtractSetConfigInputs()

	assert.Equal(t, [32]uint8{2, 1}, groupQuorums)
	assert.Equal(t, [32]uint8{0, 0}, groupParents)
	assert.Equal(t, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2"), common.HexToAddress("0x3")}, signerAddresses)
	assert.Equal(t, []uint8{0, 0, 1}, signerGroups)
}

// Test case 2: Valid configuration with only root signers
// Configuration:
// Quorum: 1
// Signers: [0x1, 0x2]
// Group signers: []
func TestExtractSetConfigInputs_OnlyRootSigners(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	config, err := NewConfig(1, signers, []Config{})
	assert.NotNil(t, config)
	assert.NoError(t, err)

	groupQuorums, groupParents, signerAddresses, signerGroups := config.ExtractSetConfigInputs()

	assert.Equal(t, [32]uint8{1, 0}, groupQuorums)
	assert.Equal(t, [32]uint8{0, 0}, groupParents)
	assert.Equal(t, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}, signerAddresses)
	assert.Equal(t, []uint8{0, 0}, signerGroups)
}

// Test case 3: Valid configuration with only groups
// Configuration:
// Quorum: 1
// Signers: []
//
//	Group signers: [{
//		 Quorum: 1
//		 Signers: [0x3]
//		 Group signers: []
//	},
//
//	{
//	  Quorum: 1
//	  Signers: [0x4]
//	  Group signers: []
//	},
//
//	{
//		 Quorum: 1
//		 Signers: [0x5]
//		 Group signers: []
//	}]
func TestExtractSetConfigInputs_OnlyGroups(t *testing.T) {
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x4")}},
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x5")}},
	}
	config, err := NewConfig(2, []common.Address{}, groupSigners)
	assert.NotNil(t, config)
	assert.NoError(t, err)

	groupQuorums, groupParents, signerAddresses, signerGroups := config.ExtractSetConfigInputs()

	assert.Equal(t, [32]uint8{2, 1, 1, 1}, groupQuorums)
	assert.Equal(t, [32]uint8{0, 0, 0, 0}, groupParents)
	assert.Equal(t, []common.Address{common.HexToAddress("0x3"), common.HexToAddress("0x4"), common.HexToAddress("0x5")}, signerAddresses)
	assert.Equal(t, []uint8{1, 2, 3}, signerGroups)
}

// Test case 4: Valid configuration with nested signers and groups
// Configuration:
// Quorum: 2
// Signers: [0x1, 0x2]
//
//		Group signers: [{
//			Quorum: 1
//			Signers: [0x3]
//			Group signers: [{
//				Quorum: 1
//				Signers: [0x4]
//				Group signers: []
//			}]
//		},
//	 {
//			Quorum: 1
//			Signers: [0x5]
//			Group signers: []
//		}]
func TestExtractSetConfigInputs_NestedSignersAndGroups(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}, GroupSigners: []Config{
			{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x4")}},
		}},
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x5")}},
	}
	config, err := NewConfig(2, signers, groupSigners)
	assert.NotNil(t, config)
	assert.NoError(t, err)

	groupQuorums, groupParents, signerAddresses, signerGroups := config.ExtractSetConfigInputs()

	assert.Equal(t, [32]uint8{2, 1, 1, 1}, groupQuorums)
	assert.Equal(t, [32]uint8{0, 0, 1, 0}, groupParents)
	assert.Equal(t, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2"), common.HexToAddress("0x3"), common.HexToAddress("0x4"), common.HexToAddress("0x5")}, signerAddresses)
	assert.Equal(t, []uint8{0, 0, 1, 2, 3}, signerGroups)
}

// Test case 5: Valid configuration with unsorted signers and groups
// Configuration:
// Quorum: 2
// Signers: [0x2, 0x1]
//
//		Group signers: [{
//			Quorum: 1
//			Signers: [0x3]
//			Group signers: [{
//				Quorum: 1
//				Signers: [0x4]
//				Group signers: []
//			}]
//		},
//	 	{
//			Quorum: 1
//			Signers: [0x5]
//			Group signers: []
//		}]
func TestExtractSetConfigInputs_UnsortedSignersAndGroups(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x2"), common.HexToAddress("0x1")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}, GroupSigners: []Config{
			{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x4")}},
		}},
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x5")}},
	}
	config, err := NewConfig(2, signers, groupSigners)
	assert.NotNil(t, config)
	assert.NoError(t, err)

	groupQuorums, groupParents, signerAddresses, signerGroups := config.ExtractSetConfigInputs()

	assert.Equal(t, [32]uint8{2, 1, 1, 1}, groupQuorums)
	assert.Equal(t, [32]uint8{0, 0, 1, 0}, groupParents)
	assert.Equal(t, []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2"), common.HexToAddress("0x3"), common.HexToAddress("0x4"), common.HexToAddress("0x5")}, signerAddresses)
	assert.Equal(t, []uint8{0, 0, 1, 2, 3}, signerGroups)
}

func TestConfigEquals_Success(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	config1, err := NewConfig(2, signers, groupSigners)
	assert.NoError(t, err)
	assert.NotNil(t, config1)

	config2, err := NewConfig(2, signers, groupSigners)
	assert.NoError(t, err)
	assert.NotNil(t, config2)

	assert.True(t, config1.Equals(config2))
}
func TestConfigEquals_Failure_MismatchingQuorum(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	config1, err := NewConfig(2, signers, groupSigners)
	assert.NoError(t, err)
	assert.NotNil(t, config1)

	config2, err := NewConfig(1, signers, groupSigners)
	assert.NoError(t, err)
	assert.NotNil(t, config2)

	assert.False(t, config1.Equals(config2))
}

func TestConfigEquals_Failure_MismatchingSigners(t *testing.T) {
	signers1 := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	signers2 := []common.Address{common.HexToAddress("0x1")}
	groupSigners := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	config1, err := NewConfig(2, signers1, groupSigners)
	assert.NotNil(t, config1)
	assert.NoError(t, err)

	config2, err := NewConfig(2, signers2, groupSigners)
	assert.NotNil(t, config2)
	assert.NoError(t, err)

	assert.False(t, config1.Equals(config2))
}

func TestConfigEquals_Failure_MismatchingGroupSigners(t *testing.T) {
	signers := []common.Address{common.HexToAddress("0x1"), common.HexToAddress("0x2")}
	groupSigners1 := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x3")}},
	}
	groupSigners2 := []Config{
		{Quorum: 1, Signers: []common.Address{common.HexToAddress("0x4")}},
	}
	config1, err := NewConfig(2, signers, groupSigners1)
	assert.NotNil(t, config1)
	assert.NoError(t, err)

	config2, err := NewConfig(2, signers, groupSigners2)
	assert.NotNil(t, config2)
	assert.NoError(t, err)

	assert.False(t, config1.Equals(config2))
}
