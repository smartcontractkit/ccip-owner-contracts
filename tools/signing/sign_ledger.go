package signing

import (
	"os"

	// NOTE MUST BE > 1.14 for this fix
	// https://github.com/ethereum/go-ethereum/pull/28945

	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

// Just run this locally to sign from the ledger.
func SignLedger(derivationPath []uint32, filePath string, proposalType proposal.ProposalType) error {
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

	// Load ledger
	ledgerhub, _ := usbwallet.NewLedgerHub()
	wallets := ledgerhub.Wallets()
	wallet := wallets[0]

	// Create executor
	executor, err := proposal.ToExecutor(make(map[mcms.ChainIdentifier]mcms.ContractDeployBackend)) // TODO: pass in a real backend
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

	// Write proposal to file
	WriteProposalToFile(proposal, os.Args[0])

	// Close wallet
	err = wallet.Close()
	if err != nil {
		return err
	}

	return nil
}
