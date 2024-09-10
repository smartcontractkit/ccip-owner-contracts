package merkle

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type MerkleTree struct {
	// Root hash of the merkle tree
	Root common.Hash

	// Layers of the merkle tree, starting from the leaves
	Layers [][]common.Hash
}

func NewMerkleTree(leaves []common.Hash) *MerkleTree {
	layers := make([][]common.Hash, 0)

	currHashes := leaves
	for len(currHashes) > 1 {
		// If the number of hashes is odd, duplicate the last hash
		if len(currHashes)%2 != 0 {
			currHashes = append(currHashes, currHashes[len(currHashes)-1])
		}

		// Append the current layer to the tree
		layers = append(layers, currHashes)

		// Calculate the parent hashes
		tempHashes := make([]common.Hash, len(currHashes)/2)
		for i := 0; i < len(currHashes); i += 2 {
			// Sort the pair of hashes before hashing
			tempHashes[i/2] = hashPair(currHashes[i], currHashes[i+1])
		}

		// Set the current hashes to the parent hashes
		currHashes = tempHashes
	}

	// Append the root hash to the tree
	return &MerkleTree{
		Root:   currHashes[0],
		Layers: layers,
	}
}

func (t *MerkleTree) GetProof(hash common.Hash) ([]common.Hash, error) {
	proof := make([]common.Hash, 0)

	targetHash := hash
	for i := 0; i < len(t.Layers); i++ {
		found := false
		for j, h := range t.Layers[i] {
			if h != targetHash {
				continue
			}

			// Get the sibling hash
			siblingIdx := j ^ 1
			siblingHash := t.Layers[i][siblingIdx]
			proof = append(proof, siblingHash)

			// Get next target hash by sorting the pair of hashes and hashing them
			targetHash = hashPair(targetHash, siblingHash)

			// Move to the next layer
			found = true
			break
		}

		if !found {
			// If the hash is not found in the current layer, it is not in the tree
			// THIS SHOULD NEVER HAPPEN
			return nil, &ErrMerkleTreeNodeNotFound{
				TargetHash: targetHash,
			}
		}
	}

	return proof, nil
}

func (t *MerkleTree) GetProofs() (map[common.Hash][]common.Hash, error) {
	proofs := make(map[common.Hash][]common.Hash)

	for _, leaf := range t.Layers[0] {
		proof, err := t.GetProof(leaf)
		if err != nil {
			// THIS SHOULD NEVER HAPPEN
			return nil, err
		}

		proofs[leaf] = proof
	}

	return proofs, nil
}

type ErrMerkleTreeNodeNotFound struct {
	TargetHash common.Hash
}

func (e *ErrMerkleTreeNodeNotFound) Error() string {
	return "merkle tree does not contain hash: " + e.TargetHash.String()
}

func hashPair(a, b common.Hash) common.Hash {
	if a.Cmp(b) < 0 {
		return efficientHash(a, b)
	} else {
		return efficientHash(b, a)
	}
}

func efficientHash(a, b common.Hash) common.Hash {
	// Create a buffer of size 64 bytes to store both hashes
	var combinedHash [64]byte

	// Copy the first hash to the first 32 bytes of the buffer
	copy(combinedHash[:32], a[:])

	// Copy the second hash to the next 32 bytes of the buffer
	copy(combinedHash[32:], b[:])

	// Compute the Keccak256 hash of the combined data
	return crypto.Keccak256Hash(combinedHash[:])
}
