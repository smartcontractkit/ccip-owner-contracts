package mcms

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/merkle"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
)

var MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP = crypto.Keccak256Hash([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP"))
var MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA = crypto.Keccak256Hash([]byte("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA"))

func calculateTransactionCounts(transactions []ChainOperation) map[ChainIdentifier]uint64 {
	txCounts := make(map[ChainIdentifier]uint64)
	for _, tx := range transactions {
		txCounts[tx.ChainIdentifier]++
	}
	return txCounts
}

func buildRootMetadatas(
	chainMetadata map[ChainIdentifier]ChainMetadata,
	txCounts map[ChainIdentifier]uint64,
	overridePreviousRoot bool,
	isSim bool,
) (map[ChainIdentifier]gethwrappers.ManyChainMultiSigRootMetadata, error) {
	rootMetadatas := make(map[ChainIdentifier]gethwrappers.ManyChainMultiSigRootMetadata)

	for chainID, metadata := range chainMetadata {
		chain, exists := chain_selectors.ChainBySelector(uint64(chainID))
		if !exists {
			return nil, &errors.ErrInvalidChainID{
				ReceivedChainID: uint64(chainID),
			}
		}

		currentTxCount, ok := txCounts[chainID]
		if !ok {
			return nil, &errors.ErrMissingChainDetails{
				ChainIdentifier: uint64(chainID),
				Parameter:       "transaction count",
			}
		}

		// Simulated chains always have block.chainid = 1337
		// So for setRoot to execute (not throw WrongChainId) we must
		// override the evmChainID to be 1337.
		if isSim {
			chain.EvmChainID = 1337
		}
		rootMetadatas[chainID] = gethwrappers.ManyChainMultiSigRootMetadata{
			ChainId:              new(big.Int).SetUint64(chain.EvmChainID),
			MultiSig:             metadata.MCMAddress,
			PreOpCount:           big.NewInt(int64(metadata.StartingOpCount)),                         // TODO: handle overflow
			PostOpCount:          big.NewInt(int64(metadata.StartingOpCount) + int64(currentTxCount)), // TODO: handle overflow
			OverridePreviousRoot: overridePreviousRoot,
		}
	}
	return rootMetadatas, nil
}

func buildOperations(
	transactions []ChainOperation,
	rootMetadatas map[ChainIdentifier]gethwrappers.ManyChainMultiSigRootMetadata,
	txCounts map[ChainIdentifier]uint64,
) (map[ChainIdentifier][]gethwrappers.ManyChainMultiSigOp, []gethwrappers.ManyChainMultiSigOp) {
	ops := make(map[ChainIdentifier][]gethwrappers.ManyChainMultiSigOp)
	chainAgnosticOps := make([]gethwrappers.ManyChainMultiSigOp, 0)
	chainIdx := make(map[ChainIdentifier]uint32, len(rootMetadatas))

	for _, tx := range transactions {
		rootMetadata := rootMetadatas[tx.ChainIdentifier]
		if _, ok := ops[tx.ChainIdentifier]; !ok {
			ops[tx.ChainIdentifier] = make([]gethwrappers.ManyChainMultiSigOp, txCounts[tx.ChainIdentifier])
			chainIdx[tx.ChainIdentifier] = 0
		}

		op := gethwrappers.ManyChainMultiSigOp{
			ChainId:  rootMetadata.ChainId,
			MultiSig: rootMetadata.MultiSig,
			Nonce:    big.NewInt(rootMetadata.PreOpCount.Int64() + int64(chainIdx[tx.ChainIdentifier])),
			To:       tx.To,
			Data:     tx.Data,
			Value:    tx.Value,
		}

		chainAgnosticOps = append(chainAgnosticOps, op)
		ops[tx.ChainIdentifier][chainIdx[tx.ChainIdentifier]] = op
		chainIdx[tx.ChainIdentifier]++
	}

	return ops, chainAgnosticOps
}

func sortedChainIdentifiers(chainMetadata map[ChainIdentifier]ChainMetadata) []ChainIdentifier {
	chainIdentifiers := make([]ChainIdentifier, 0, len(chainMetadata))
	for chainID := range chainMetadata {
		chainIdentifiers = append(chainIdentifiers, chainID)
	}
	sort.Slice(chainIdentifiers, func(i, j int) bool { return chainIdentifiers[i] < chainIdentifiers[j] })
	return chainIdentifiers
}

func buildMerkleTree(
	chainIdentifiers []ChainIdentifier,
	rootMetadatas map[ChainIdentifier]gethwrappers.ManyChainMultiSigRootMetadata,
	ops map[ChainIdentifier][]gethwrappers.ManyChainMultiSigOp,
) (*merkle.MerkleTree, error) {
	hashLeaves := make([]common.Hash, 0)

	for _, chainID := range chainIdentifiers {
		encodedRootMetadata, err := metadataEncoder(rootMetadatas[chainID])
		if err != nil {
			return nil, err
		}
		hashLeaves = append(hashLeaves, encodedRootMetadata)

		for _, op := range ops[chainID] {
			encodedOp, err := txEncoder(op)
			if err != nil {
				return nil, err
			}
			hashLeaves = append(hashLeaves, encodedOp)
		}
	}

	// sort the hashes and sort the pairs
	sort.Slice(hashLeaves, func(i, j int) bool {
		return hashLeaves[i].String() < hashLeaves[j].String()
	})

	return merkle.NewMerkleTree(hashLeaves), nil
}

func metadataEncoder(rootMetadata gethwrappers.ManyChainMultiSigRootMetadata) (common.Hash, error) {
	abi := `[{"type":"bytes32"},{"type":"tuple","components":[{"name":"chainId","type":"uint256"},{"name":"multiSig","type":"address"},{"name":"preOpCount","type":"uint40"},{"name":"postOpCount","type":"uint40"},{"name":"overridePreviousRoot","type":"bool"}]}]`
	encoded, err := ABIEncode(abi, MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA, rootMetadata)
	if err != nil {
		return common.Hash{}, err
	}

	return crypto.Keccak256Hash(encoded), nil
}

func txEncoder(op gethwrappers.ManyChainMultiSigOp) (common.Hash, error) {
	abi := `[{"type":"bytes32"},{"type":"tuple","components":[{"name":"chainId","type":"uint256"},{"name":"multiSig","type":"address"},{"name":"nonce","type":"uint40"},{"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"data","type":"bytes"}]}]`
	encoded, err := ABIEncode(abi, MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, op)
	if err != nil {
		return common.Hash{}, err
	}

	return crypto.Keccak256Hash(encoded), nil
}
