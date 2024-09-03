package timelockproposal

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	owner "github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/mcmsproposal"
)

var ZERO_HASH = common.Hash{}

type MCMSWithTimelockChainMetadata struct {
	mcmsproposal.ChainMetadata
	TimelockAddress common.Address `json:"timelockAddress"`
}

type TimelockOperation string

const (
	Schedule TimelockOperation = "schedule"
	Cancel   TimelockOperation = "cancel"
	Bypass   TimelockOperation = "bypass"
)

type MCMSWithTimelockProposal struct {
	mcmsproposal.Proposal

	Operation TimelockOperation `json:"operation"` // Always 'schedule', 'cancel', or 'bypass'

	// TODO: this should be configurable as a human-readable string
	// i.e. 1d, 1w, 1m, 1y
	MinDelay string `json:"minDelay"`

	// Overridden: Map of chain identifier to chain metadata
	ChainMetadata map[string]MCMSWithTimelockChainMetadata `json:"chainMetadata"`

	// Overridden: Operations to be executed after wrapping in a timelock
	Transactions []BatchChainOperation `json:"transactions"`
}

func (m *MCMSWithTimelockProposal) Validate() error {
	if err := m.Proposal.Validate(); err != nil {
		return err
	}

	switch m.Operation {
	case Schedule, Cancel, Bypass:
		break
	default:
		return &errors.ErrInvalidTimelockOperation{
			ReceivedTimelockOperation: string(m.Operation),
		}
	}

	_, err := time.ParseDuration(m.MinDelay)
	if err != nil {
		return &errors.ErrInvalidMinDelay{
			ReceivedMinDelay: m.MinDelay,
		}
	}

	return nil
}

func (m *MCMSWithTimelockProposal) ToMCMSOnlyProposal() (mcmsproposal.Proposal, error) {
	mcmOnly := m.Proposal

	// Start predecessor map with all chains pointing to the zero hash
	predecessorMap := make(map[string]common.Hash)
	for chain := range m.ChainMetadata {
		predecessorMap[chain] = ZERO_HASH
	}

	// Convert chain metadata
	mcmOnly.ChainMetadata = make(map[string]mcmsproposal.ChainMetadata)
	for chain, metadata := range m.ChainMetadata {
		mcmOnly.ChainMetadata[chain] = mcmsproposal.ChainMetadata{
			NonceOffset: metadata.NonceOffset,
			MCMAddress:  metadata.MCMAddress,
		}
	}

	// Convert transactions into timelock wrapped transactions
	for _, t := range m.Transactions {
		calls := make([]owner.RBACTimelockCall, 0)
		tags := make([]string, 0)
		for _, op := range t.Batch {
			calls = append(calls, owner.RBACTimelockCall{
				Target: op.To,
				Data:   op.Data,
				Value:  op.Value,
			})
			tags = append(tags, op.Tags...)
		}
		predecessor := predecessorMap[t.ChainIdentifier]
		salt := ZERO_HASH
		delay, _ := time.ParseDuration(m.MinDelay)

		abi, err := owner.RBACTimelockMetaData.GetAbi()
		if err != nil {
			return mcmsproposal.Proposal{}, err
		}
		data, err := abi.Pack("scheduleBatch", calls, predecessor, salt, big.NewInt(int64(delay.Seconds())))
		if err != nil {
			return mcmsproposal.Proposal{}, err
		}

		mcmOnly.Transactions = append(mcmOnly.Transactions, mcmsproposal.ChainOperation{
			ChainIdentifier: t.ChainIdentifier,
			Operation: mcmsproposal.Operation{
				To:           m.ChainMetadata[t.ChainIdentifier].TimelockAddress,
				Data:         data,
				Value:        big.NewInt(0), // TODO: is this right?
				ContractType: "RBACTimelock",
				Tags:         tags,
			},
		})

		predecessorMap[t.ChainIdentifier], err = hashOperationBatch(calls, predecessor, salt)
		if err != nil {
			return mcmsproposal.Proposal{}, err
		}
	}

	return mcmOnly, nil
}

// hashOperationBatch replicates the hash calculation from Solidity
// TODO: see if there's an easier way to do this using the gethwrappers
func hashOperationBatch(calls []owner.RBACTimelockCall, predecessor, salt [32]byte) ([32]byte, error) {
	// Encode the calls using RLP encoding
	encodedCalls, err := rlp.EncodeToBytes(calls)
	if err != nil {
		return [32]byte{}, err
	}

	// Encode the entire data (calls, predecessor, salt) using ABI encoding
	encoded := crypto.Keccak256(
		append(encodedCalls, append(predecessor[:], salt[:]...)...),
	)

	// Return the hash as a [32]byte array
	return [32]byte(crypto.Keccak256Hash(encoded).Bytes()), nil
}
