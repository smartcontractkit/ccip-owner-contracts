package signing

import (
	"crypto/ecdsa"
	"os"

	// NOTE MUST BE > 1.14 for this fix
	// https://github.com/ethereum/go-ethereum/pull/28945

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

// Just run this locally to sign from the ledger.
func SignPlainKey(privateKey *ecdsa.PrivateKey, filePath string, proposalType proposal.ProposalType) error {
	// Load file
	proposal, err := LoadProposal(proposalType, filePath)
	if err != nil {
		return err
	}

	// Validate proposal
	err = proposal.Validate()
	if err != nil {
		return err
	}

	// Create executor
	executor, err := proposal.ToExecutor(make(map[mcms.ChainIdentifier]mcms.ContractDeployBackend)) // TODO: pass in a real backend
	if err != nil {
		return err
	}

	// Get the signing hash
	payload, err := executor.SigningHash()
	if err != nil {
		return err
	}

	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), privateKey)
	if err != nil {
		return err
	}

	// Unmarshal signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	if err != nil {
		return err
	}

	// Add signature to proposal
	proposal.AddSignature(sigObj)

	// Write proposal to file
	WriteProposalToFile(proposal, os.Args[0])
	return nil
}
