package configwrappers

import (
	"math/big"
	"sort"

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

func (c *Config) ToRawConfig() gethwrappers.ManyChainMultiSigConfig {
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
	}
}

func (c *Config) Equals(other *Config) bool {
	if c.Quorum != other.Quorum {
		return false
	}

	if len(c.Signers) != len(other.Signers) {
		return false
	}

	// Compare signers (order doesn't matter)
	if !unorderedArrayEquals(c.Signers, other.Signers) {
		return false
	}

	if len(c.GroupSigners) != len(other.GroupSigners) {
		return false
	}

	// Compare all group signers in first exist in second (order doesn't matter)
	for i := range c.GroupSigners {
		found := false
		for j := range other.GroupSigners {
			if c.GroupSigners[i].Equals(&other.GroupSigners[j]) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	// Compare all group signers in second exist in first (order doesn't matter)
	for i := range other.GroupSigners {
		found := false
		for j := range c.GroupSigners {
			if other.GroupSigners[i].Equals(&c.GroupSigners[j]) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
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

	// Combine SignerAddresses and SignerGroups into a slice of Signer structs
	signerObjs := make([]gethwrappers.ManyChainMultiSigSigner, len(signers))
	for i := range signers {
		signerObjs[i] = gethwrappers.ManyChainMultiSigSigner{
			Addr:  signers[i],
			Group: signerGroups[i],
		}
	}

	// Sort signers by their addresses in ascending order
	sort.Slice(signerObjs, func(i, j int) bool {
		addressA := new(big.Int).SetBytes(signerObjs[i].Addr.Bytes())
		addressB := new(big.Int).SetBytes(signerObjs[j].Addr.Bytes())
		return addressA.Cmp(addressB) < 0
	})

	// Extract the ordered addresses and groups after sorting
	orderedSignerAddresses := make([]common.Address, len(signers))
	orderedSignerGroups := make([]uint8, len(signers))
	for i, signer := range signerObjs {
		orderedSignerAddresses[i] = signer.Addr
		orderedSignerGroups[i] = signer.Group
	}

	return [32]uint8(groupQuorums), [32]uint8(groupParents), orderedSignerAddresses, orderedSignerGroups
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

func unorderedArrayEquals[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[T]struct{})
	bMap := make(map[T]struct{})

	for _, i := range a {
		aMap[i] = struct{}{}
	}

	for _, i := range b {
		bMap[i] = struct{}{}
	}

	for _, i := range a {
		if _, ok := bMap[i]; !ok {
			return false
		}
	}

	for _, i := range b {
		if _, ok := aMap[i]; !ok {
			return false
		}
	}

	return true
}
