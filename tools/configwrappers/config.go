package configwrappers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

// Config is a struct that holds all the configuration for the owner contracts
type Config struct {
	Quorum uint8 `json:"quorum"`

	// TODO: how should this change as we expand to other non-EVM chains?
	Signers      []common.Address `json:"signers"`
	GroupSigners []Config         `json:"groupSigners"`
}

func NewConfig(quorum uint8, signers []common.Address, groupSigners []Config) *Config {
	return &Config{
		Quorum:       quorum,
		Signers:      signers,
		GroupSigners: groupSigners,
	}
}

func NewConfigFromRaw(rawConfig gethwrappers.ManyChainMultiSigConfig) *Config {
	groupToSigners := make([][]common.Address, len(rawConfig.GroupQuorums))
	for _, signer := range rawConfig.Signers {
		groupToSigners[signer.Group] = append(groupToSigners[signer.Group], signer.Addr)
	}

	groups := make([]Config, len(rawConfig.GroupQuorums))
	for i, quorum := range rawConfig.GroupQuorums {
		signers := groupToSigners[i]
		if signers == nil {
			signers = []common.Address{}
		}

		groups[i] = Config{
			Signers:      signers,
			GroupSigners: []Config{},
			Quorum:       quorum,
		}
	}

	for i, parent := range rawConfig.GroupParents {
		if i > 0 && groups[i].Quorum > 0 {
			groups[parent].GroupSigners = append(groups[parent].GroupSigners, groups[i])
		}
	}

	return &groups[0]
}

func (c *Config) Validate() error {
	if c.Quorum == 0 {
		return &errors.ErrInvalidMCMSConfig{
			Reason: "Quorum must be greater than 0",
		}
	}

	if len(c.Signers) == 0 && len(c.GroupSigners) == 0 {
		return &errors.ErrInvalidMCMSConfig{
			Reason: "Config must have at least one signer or group",
		}
	}

	if (len(c.Signers) + len(c.GroupSigners)) < int(c.Quorum) {
		return &errors.ErrInvalidMCMSConfig{
			Reason: "Quorum must be less than or equal to the number of signers and groups",
		}
	}

	for _, groupSigner := range c.GroupSigners {
		if err := groupSigner.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) ToRawConfig() (gethwrappers.ManyChainMultiSigConfig, error) {
	groupQuorums, groupParents, signerAddresses, signerGroups := c.ExtractSetConfigInputs()

	// convert to gethwrappers types
	signers := make([]gethwrappers.ManyChainMultiSigSigner, len(signerAddresses))
	for i, signer := range signerAddresses {
		signers[i] = gethwrappers.ManyChainMultiSigSigner{
			Addr:  signer,
			Group: signerGroups[i],
			Index: uint8(i),
		}
	}

	return gethwrappers.ManyChainMultiSigConfig{
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
		Signers:      signers,
	}, nil
}

func (c *Config) ExtractSetConfigInputs() ([32]uint8, [32]uint8, []common.Address, []uint8) {
	var groupQuorums, groupParents, signerGroups []uint8 = []uint8{}, []uint8{}, []uint8{}
	var signers []common.Address = []common.Address{}

	extractGroupsAndSigners(c, 0, &groupQuorums, &groupParents, &signers, &signerGroups)

	// fill the rest of the arrays with 0s
	for i := len(groupQuorums); i < 32; i++ {
		groupQuorums = append(groupQuorums, 0)
		groupParents = append(groupParents, 0)
	}

	return [32]uint8(groupQuorums), [32]uint8(groupParents), signers, signerGroups
}

func extractGroupsAndSigners(group *Config, parentIdx int, groupQuorums *[]uint8, groupParents *[]uint8, signers *[]common.Address, signerGroups *[]uint8) {
	// Append the group's quorum and parent index to the respective slices
	*groupQuorums = append(*groupQuorums, group.Quorum)
	*groupParents = append(*groupParents, uint8(parentIdx))

	// Assign the current group index
	currentGroupIdx := len(*groupQuorums) - 1

	// For each string signer, append the signer and its group index
	for _, signer := range group.Signers {
		*signers = append(*signers, signer)
		*signerGroups = append(*signerGroups, uint8(currentGroupIdx))
	}

	// Recursively handle the nested multisig groups
	for _, groupSigner := range group.GroupSigners {
		extractGroupsAndSigners(&groupSigner, currentGroupIdx, groupQuorums, groupParents, signers, signerGroups)
	}
}
