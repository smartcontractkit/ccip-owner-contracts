package executable

import (
	"encoding/binary"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

type Executor struct {
	Proposal         *ExecutableMCMSProposal
	Tree             *MerkleTree
	RootMetadatas    map[string]gethwrappers.ManyChainMultiSigRootMetadata
	Operations       map[string][]gethwrappers.ManyChainMultiSigOp
	ChainAgnosticOps []gethwrappers.ManyChainMultiSigOp
	Caller           *Caller
}

func NewProposalExecutor(proposal *ExecutableMCMSProposal, clients map[string]ContractDeployBackend) (*Executor, error) {
	txCounts := calculateTransactionCounts(proposal.Transactions)

	caller, err := NewCaller(mapMCMAddresses(proposal.ChainMetadata), clients)
	if err != nil {
		return nil, err
	}

	currentOpCounts, err := caller.GetCurrentOpCounts()
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

	return &Executor{
		Proposal:         proposal,
		Tree:             tree,
		RootMetadatas:    rootMetadatas,
		Operations:       ops,
		ChainAgnosticOps: chainAgnosticOps,
		Caller:           caller,
	}, err
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

func (e *Executor) ValidateSignatures() error {
	hash, err := e.SigningHash()
	if err != nil {
		return err
	}

	recoveredSigners := make([]common.Address, len(e.Proposal.Signatures))
	for i, sig := range e.Proposal.Signatures {
		recoveredAddr, err := sig.Recover(hash)
		if err != nil {
			return err
		}
		recoveredSigners[i] = recoveredAddr
	}

	configs, err := e.Caller.GetConfigs()
	if err != nil {
		return err
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
				return &errors.ErrInvalidSignature{
					ChainIdentifier:  chain,
					RecoveredAddress: signer,
				}
			}
		}
	}

	return nil
}

func (e *Executor) SetRootOnChain(auth *bind.TransactOpts, chain string) (*types.Transaction, error) {
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

	return e.Caller.Callers[chain].SetRoot(
		auth,
		[32]byte(e.Tree.Root.Bytes()),
		e.Proposal.ValidUntil, metadata,
		mapHashes(proof),
		mapSignatures(sortedSignatures),
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

	return e.Caller.Callers[chain].Execute(
		auth,
		mcmOperation,
		mapHashes(proof),
	)
}
