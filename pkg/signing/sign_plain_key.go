package signing

import (
	"crypto/ecdsa"

	// NOTE MUST BE > 1.14 for this fix
	// https://github.com/ethereum/go-ethereum/pull/28945

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
)

// Just run this locally to sign from the ledger.
func SignPlainKey(privateKey *ecdsa.PrivateKey, proposal proposal.Proposal) error {
	// Validate proposal
	err := proposal.Validate()
	if err != nil {
		return err
	}

	executor, err := proposal.ToExecutor(false) // TODO: pass in a real backend
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
	return nil
}
