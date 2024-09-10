package timelock

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	owner "github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

var ZERO_HASH = common.Hash{}

type TimelockOperation string

const (
	Schedule TimelockOperation = "schedule"
	Cancel   TimelockOperation = "cancel"
	Bypass   TimelockOperation = "bypass"
)

type MCMSWithTimelockProposal struct {
	mcms.Proposal

	Operation TimelockOperation `json:"operation"` // Always 'schedule', 'cancel', or 'bypass'

	// i.e. 1d, 1w, 1m, 1y
	MinDelay string `json:"minDelay"`

	TimelockAddresses map[mcms.ChainIdentifier]common.Address `json:"timelockAddresses"`

	// Overridden: Operations to be executed after wrapping in a timelock
	Transactions []BatchChainOperation `json:"transactions"`
}

func NewMCMSWithTimelockProposal(
	version string,
	validUntil uint32,
	signatures []mcms.Signature,
	overridePreviousRoot bool,
	chainMetadata map[mcms.ChainIdentifier]mcms.ChainMetadata,
	timelockAddresses map[mcms.ChainIdentifier]common.Address,
	description string,
	transactions []BatchChainOperation,
	operation TimelockOperation,
	minDelay string,
) (*MCMSWithTimelockProposal, error) {
	proposal := MCMSWithTimelockProposal{
		Proposal: mcms.Proposal{
			Version:              version,
			ValidUntil:           validUntil,
			Signatures:           signatures,
			OverridePreviousRoot: overridePreviousRoot,
			Description:          description,
			ChainMetadata:        chainMetadata,
		},
		Operation:         operation,
		MinDelay:          minDelay,
		TimelockAddresses: timelockAddresses,
		Transactions:      transactions,
	}

	err := proposal.Validate()
	if err != nil {
		return nil, err
	}

	return &proposal, nil
}

func (m *MCMSWithTimelockProposal) Validate() error {
	if m.Version == "" {
		return &errors.ErrInvalidVersion{
			ReceivedVersion: m.Version,
		}
	}

	// Get the current Unix timestamp as an int64
	currentTime := time.Now().Unix()

	if m.ValidUntil <= uint32(currentTime) {
		// ValidUntil is a Unix timestamp, so it should be greater than the current time
		return &errors.ErrInvalidValidUntil{
			ReceivedValidUntil: m.ValidUntil,
		}
	}

	if len(m.ChainMetadata) == 0 {
		return &errors.ErrNoChainMetadata{}
	}

	if len(m.Transactions) == 0 {
		return &errors.ErrNoTransactions{}
	}

	if m.Description == "" {
		return &errors.ErrInvalidDescription{
			ReceivedDescription: m.Description,
		}
	}

	// Validate all chains in transactions have an entry in chain metadata
	for _, t := range m.Transactions {
		if _, ok := m.ChainMetadata[t.ChainIdentifier]; !ok {
			return &errors.ErrMissingChainDetails{
				ChainIdentifier: uint64(t.ChainIdentifier),
				Parameter:       "chain metadata",
			}
		}
	}

	switch m.Operation {
	case Schedule, Cancel, Bypass:
		break
	default:
		return &errors.ErrInvalidTimelockOperation{
			ReceivedTimelockOperation: string(m.Operation),
		}
	}

	// Validate the delay is a valid duration but is only required
	// for Schedule operations
	if m.Operation == Schedule {
		if _, err := time.ParseDuration(m.MinDelay); err != nil {
			return err
		}
	}

	return nil
}

func (m *MCMSWithTimelockProposal) ToExecutor(clients map[mcms.ChainIdentifier]mcms.ContractDeployBackend) (*mcms.Executor, error) {
	// Convert the proposal to an MCMS only proposal
	mcmOnly, err := m.toMCMSOnlyProposal()
	if err != nil {
		return nil, err
	}

	return mcmOnly.ToExecutor(clients)
}

func (m *MCMSWithTimelockProposal) toMCMSOnlyProposal() (mcms.Proposal, error) {
	mcmOnly := m.Proposal

	// Start predecessor map with all chains pointing to the zero hash
	predecessorMap := make(map[mcms.ChainIdentifier]common.Hash)
	for chain := range m.ChainMetadata {
		predecessorMap[chain] = ZERO_HASH
	}

	// Convert chain metadata
	mcmOnly.ChainMetadata = make(map[mcms.ChainIdentifier]mcms.ChainMetadata)
	for chain, metadata := range m.ChainMetadata {
		mcmOnly.ChainMetadata[chain] = mcms.ChainMetadata{
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
			return mcms.Proposal{}, err
		}

		operationId, err := hashOperationBatch(calls, predecessor, salt)
		if err != nil {
			return mcms.Proposal{}, err
		}

		// Encode the data based on the operation
		var data []byte
		switch m.Operation {
		case Schedule:
			data, err = abi.Pack("scheduleBatch", calls, predecessor, salt, big.NewInt(int64(delay.Seconds())))
			if err != nil {
				return mcms.Proposal{}, err
			}
		case Cancel:
			data, err = abi.Pack("cancel", operationId)
			if err != nil {
				return mcms.Proposal{}, err
			}
		case Bypass:
			data, err = abi.Pack("bypasserExecuteBatch", calls)
			if err != nil {
				return mcms.Proposal{}, err
			}
		default:
			return mcms.Proposal{}, &errors.ErrInvalidTimelockOperation{
				ReceivedTimelockOperation: string(m.Operation),
			}
		}

		mcmOnly.Transactions = append(mcmOnly.Transactions, mcms.ChainOperation{
			ChainIdentifier: t.ChainIdentifier,
			Operation: mcms.Operation{
				To:           m.TimelockAddresses[t.ChainIdentifier],
				Data:         data,
				Value:        big.NewInt(0), // TODO: is this right?
				ContractType: "RBACTimelock",
				Tags:         tags,
			},
		})

		predecessorMap[t.ChainIdentifier] = operationId
	}

	return mcmOnly, nil
}

// hashOperationBatch replicates the hash calculation from Solidity
// TODO: see if there's an easier way to do this using the gethwrappers
func hashOperationBatch(calls []owner.RBACTimelockCall, predecessor, salt [32]byte) (common.Hash, error) {
	const abi = `[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"struct Call[]","name":"calls","type":"tuple[]"},{"internalType":"bytes32","name":"predecessor","type":"bytes32"},{"internalType":"bytes32","name":"salt","type":"bytes32"}]`
	encoded, err := mcms.ABIEncode(abi, calls, predecessor, salt)
	if err != nil {
		return common.Hash{}, err
	}

	// Return the hash as a [32]byte array
	return crypto.Keccak256Hash(encoded), nil
}
