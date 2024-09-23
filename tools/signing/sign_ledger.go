package signing

import (
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

// Just run this locally to sign from the ledger.
func SignLedger(derivationPath []uint32, proposal proposal.Proposal) error {
	// Validate proposal
	err := proposal.Validate()
	if err != nil {
		return err
	}

	// Load ledger
	ledgerhub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return err
	}

	// Get the first wallet
	wallets := ledgerhub.Wallets()
	wallet := wallets[0]

	// Create executor
	executor, err := proposal.ToExecutor(false) // TODO: pass in a real backend
	if err != nil {
		return err
	}

	// Open the ledger.
	_ = wallet.Open("")

	// Load account.
	// BIP44 derivation path used in ledger.
	// Could pass this in as an argument as well.
	account, err := wallet.Derive(derivationPath, true)
	if err != nil {
		return err
	}

	// Get the signing hash
	payload, err := executor.SigningHash()
	if err != nil {
		return err
	}

	// Sign the payload
	sig, err := wallet.SignData(account, "", payload.Bytes())
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

	// Close wallet
	err = wallet.Close()
	if err != nil {
		return err
	}

	return nil
}
