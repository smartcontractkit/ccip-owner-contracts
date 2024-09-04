package mcms

import (
	"encoding/binary"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/configwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/merkle"
)

type Executor struct {
	Proposal         *Proposal
	Tree             *merkle.MerkleTree
	RootMetadatas    map[ChainIdentifier]gethwrappers.ManyChainMultiSigRootMetadata
	Operations       map[ChainIdentifier][]gethwrappers.ManyChainMultiSigOp
	ChainAgnosticOps []gethwrappers.ManyChainMultiSigOp
	Callers          map[ChainIdentifier]*gethwrappers.ManyChainMultiSig
}

func NewProposalExecutor(proposal *Proposal, clients map[ChainIdentifier]ContractDeployBackend) (*Executor, error) {
	txCounts := calculateTransactionCounts(proposal.Transactions)
	mcms := transformMCMAddresses(proposal.ChainMetadata)

	mcmsWrappers := make(map[ChainIdentifier]*gethwrappers.ManyChainMultiSig)
	for chain, mcmAddress := range mcms {
		client, ok := clients[chain]
		if !ok {
			return nil, &errors.ErrMissingChainClient{
				ChainIdentifier: uint64(chain),
			}
		}

		mcms, err := gethwrappers.NewManyChainMultiSig(mcmAddress, client)
		if err != nil {
			return nil, err
		}

		mcmsWrappers[chain] = mcms
	}

	executor := &Executor{
		Proposal: proposal,
		Callers:  mcmsWrappers,
	}

	currentOpCounts, err := executor.GetCurrentOpCounts()
	if err != nil {
		return nil, err
	}

	rootMetadatas, err := buildRootMetadatas(proposal.ChainMetadata, txCounts, currentOpCounts, proposal.OverridePreviousRoot)
	if err != nil {
		return nil, err
	}

	ops, chainAgnosticOps := buildOperations(proposal.Transactions, rootMetadatas, txCounts)
	chainIdentifiers := sortedChainIdentifiers(proposal.ChainMetadata)
	tree, err := buildMerkleTree(chainIdentifiers, rootMetadatas, ops)

	executor.Tree = tree
	executor.RootMetadatas = rootMetadatas
	executor.Operations = ops
	executor.ChainAgnosticOps = chainAgnosticOps

	return executor, err
}

func (e *Executor) SigningHash() (common.Hash, error) {
	// Convert validUntil to [32]byte
	var validUntilBytes [32]byte
	binary.BigEndian.PutUint32(validUntilBytes[28:], e.Proposal.ValidUntil) // Place the uint32 in the last 4 bytes

	hashToSign := crypto.Keccak256Hash(e.Tree.Root.Bytes(), validUntilBytes[:])
	return toEthSignedMessageHash(hashToSign), nil
}

func toEthSignedMessageHash(messageHash common.Hash) common.Hash {
	// Add the Ethereum signed message prefix
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	data := append(prefix, messageHash.Bytes()...)

	// Hash the prefixed message
	return crypto.Keccak256Hash(data)
}

func (e *Executor) ValidateMCMSConfigs() error {
	configs, err := e.GetConfigs()
	if err != nil {
		return err
	}

	wrappedConfigs, err := transformMCMSConfigs(configs)
	if err != nil {
		return err
	}

	// Validate that all configs are equivalent
	sortedChains := sortedChainIdentifiers(e.Proposal.ChainMetadata)
	for i, chain := range sortedChains {
		if i == 0 {
			continue
		}

		if !wrappedConfigs[chain].Equals(wrappedConfigs[sortedChains[i-1]]) {
			return &errors.ErrInconsistentConfigs{
				ChainIdentifierA: uint64(chain),
				ChainIdentifierB: uint64(sortedChains[i-1]),
			}
		}
	}

	return nil
}

func (m *Executor) GetCurrentOpCounts() (map[ChainIdentifier]big.Int, error) {
	opCounts := make(map[ChainIdentifier]big.Int)

	for chain, wrapper := range m.Callers {
		opCount, err := wrapper.GetOpCount(&bind.CallOpts{})
		if err != nil {
			return nil, err
		}

		opCounts[chain] = *opCount
	}

	return opCounts, nil
}

func (m *Executor) GetConfigs() (map[ChainIdentifier]gethwrappers.ManyChainMultiSigConfig, error) {
	configs := make(map[ChainIdentifier]gethwrappers.ManyChainMultiSigConfig)

	for chain, wrapper := range m.Callers {
		config, err := wrapper.GetConfig(&bind.CallOpts{})
		if err != nil {
			return nil, err
		}

		configs[chain] = config
	}

	return configs, nil
}

func (e *Executor) ValidateSignatures() (bool, error) {
	hash, err := e.SigningHash()
	if err != nil {
		return false, err
	}

	recoveredSigners := make([]common.Address, len(e.Proposal.Signatures))
	for i, sig := range e.Proposal.Signatures {
		recoveredAddr, err := sig.Recover(hash)
		if err != nil {
			return false, err
		}
		recoveredSigners[i] = recoveredAddr
	}

	configs, err := e.GetConfigs()
	if err != nil {
		return false, err
	}

	// Validate that all signers are valid
	for chain, config := range configs {
		for _, signer := range recoveredSigners {
			found := false
			for _, mcmsSigner := range config.Signers {
				if mcmsSigner.Addr == signer {
					found = true
					break
				}
			}

			if !found {
				return false, &errors.ErrInvalidSignature{
					ChainIdentifier:  uint64(chain),
					RecoveredAddress: signer,
				}
			}
		}
	}

	// Validate if the quorum is met
	wrappedConfigs, err := transformMCMSConfigs(configs)
	if err != nil {
		return false, err
	}

	for chain, config := range wrappedConfigs {
		if !isReadyToSetRoot(*config, recoveredSigners) {
			return false, &errors.ErrQuorumNotMet{
				ChainIdentifier: uint64(chain),
			}
		}
	}

	return true, nil
}

func isReadyToSetRoot(rootGroup configwrappers.Config, recoveredSigners []common.Address) bool {
	return isGroupAtConsensus(rootGroup, recoveredSigners)
}

func isGroupAtConsensus(group configwrappers.Config, recoveredSigners []common.Address) bool {
	signerApprovalsInGroup := 0
	for _, signer := range group.Signers {
		for _, recoveredSigner := range recoveredSigners {
			if signer == recoveredSigner {
				signerApprovalsInGroup++
				break
			}
		}
	}

	groupApprovals := 0
	for _, groupSigner := range group.GroupSigners {
		if isGroupAtConsensus(groupSigner, recoveredSigners) {
			groupApprovals++
		}
	}

	return (signerApprovalsInGroup + groupApprovals) >= int(group.Quorum)
}

func (e *Executor) SetRootOnChain(auth *bind.TransactOpts, chain ChainIdentifier) (*types.Transaction, error) {
	metadata := e.RootMetadatas[chain]

	encodedMetadata, err := metadataEncoder(metadata)
	if err != nil {
		return nil, err
	}

	proof, err := e.Tree.GetProof(encodedMetadata)
	if err != nil {
		return nil, err
	}

	hash, err := e.SigningHash()
	if err != nil {
		return nil, err
	}

	// Sort signatures by recovered address
	sortedSignatures := e.Proposal.Signatures
	sort.Slice(sortedSignatures, func(i, j int) bool {
		recoveredSignerA, _ := sortedSignatures[i].Recover(hash)
		recoveredSignerB, _ := sortedSignatures[j].Recover(hash)
		return recoveredSignerA.Cmp(recoveredSignerB) < 0
	})

	return e.Callers[chain].SetRoot(
		auth,
		[32]byte(e.Tree.Root.Bytes()),
		e.Proposal.ValidUntil, metadata,
		transformHashes(proof),
		transformSignatures(sortedSignatures),
	)
}

func (e *Executor) ExecuteOnChain(auth *bind.TransactOpts, idx int) (*types.Transaction, error) {
	operation := e.Proposal.Transactions[idx]
	chain := operation.ChainIdentifier

	mcmOperation := e.ChainAgnosticOps[idx]
	hash, err := txEncoder(mcmOperation)
	if err != nil {
		return nil, err
	}

	proof, err := e.Tree.GetProof(hash)
	if err != nil {
		return nil, err
	}

	return e.Callers[chain].Execute(
		auth,
		mcmOperation,
		transformHashes(proof),
	)
}
