package executable

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

func calculateTransactionCounts(transactions []ChainOperation) map[string]uint64 {
	txCounts := make(map[string]uint64)
	for _, tx := range transactions {
		txCounts[tx.ChainIdentifier]++
	}
	return txCounts
}

func buildRootMetadatas(chainMetadata map[string]ExecutableMCMSChainMetadata, txCounts map[string]uint64, overridePreviousRoot bool) (map[string]gethwrappers.ManyChainMultiSigRootMetadata, error) {
	rootMetadatas := make(map[string]gethwrappers.ManyChainMultiSigRootMetadata)
	currentNonce := int64(0) // TODO: fetch this from the chain

	for chainID, metadata := range chainMetadata {
		stringChainId, err := big.NewInt(0).SetString(chainID, 10)
		if !err {
			return nil, &errors.ErrInvalidChainID{
				ReceivedChainID: chainID,
			}
		}

		rootMetadatas[chainID] = gethwrappers.ManyChainMultiSigRootMetadata{
			ChainId:              stringChainId,
			MultiSig:             common.HexToAddress(metadata.MCMAddress),
			PreOpCount:           big.NewInt(currentNonce + int64(metadata.NonceOffset)),
			PostOpCount:          big.NewInt(currentNonce + int64(metadata.NonceOffset) + int64(txCounts[chainID])),
			OverridePreviousRoot: overridePreviousRoot,
		}
	}
	return rootMetadatas, nil
}

func buildOperations(transactions []ChainOperation, rootMetadatas map[string]gethwrappers.ManyChainMultiSigRootMetadata, txCounts map[string]uint64) (map[string][]gethwrappers.ManyChainMultiSigOp, error) {
	ops := make(map[string][]gethwrappers.ManyChainMultiSigOp)
	chainIdx := make(map[string]uint32, len(rootMetadatas))

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
			To:       common.HexToAddress(tx.To),
			Data:     common.FromHex(tx.Data),
			Value:    big.NewInt(int64(tx.Value)),
		}

		ops[tx.ChainIdentifier][chainIdx[tx.ChainIdentifier]] = op
		chainIdx[tx.ChainIdentifier]++
	}

	return ops, nil
}

func sortedChainIdentifiers(chainMetadata map[string]ExecutableMCMSChainMetadata) []string {
	chainIdentifiers := make([]string, 0, len(chainMetadata))
	for chainID := range chainMetadata {
		chainIdentifiers = append(chainIdentifiers, chainID)
	}
	sort.Strings(chainIdentifiers)
	return chainIdentifiers
}

func buildMerkleTree(chainIdentifiers []string, ops map[string][]gethwrappers.ManyChainMultiSigOp) (map[string]string, string, error) {
	tree := make(map[string]string)
	currentHashes := make([]string, 0)

	for _, chainID := range chainIdentifiers {
		for _, op := range ops[chainID] {
			encoded, err := rlp.EncodeToBytes(op)
			if err != nil {
				return nil, "", err
			}

			hash := common.Bytes2Hex(crypto.Keccak256(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, encoded))
			currentHashes = append(currentHashes, hash)
		}
	}

	if len(currentHashes)%2 != 0 {
		currentHashes = append(currentHashes, currentHashes[len(currentHashes)-1])
	}

	for len(currentHashes) > 1 {
		tempCurrentHashes := make([]string, 0)

		for i := 0; i < len(currentHashes); i += 2 {
			parentHash := common.Bytes2Hex(crypto.Keccak256(common.FromHex(currentHashes[i]), common.FromHex(currentHashes[i+1])))
			tree[currentHashes[i]] = parentHash
			tree[currentHashes[i+1]] = parentHash
			tempCurrentHashes = append(tempCurrentHashes, parentHash)
		}

		currentHashes = tempCurrentHashes
	}

	root := currentHashes[0]
	return tree, root, nil
}
