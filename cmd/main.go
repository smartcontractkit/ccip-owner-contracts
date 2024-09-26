package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/spf13/cobra"
)

var (
	proposalPath string
	proposalType proposal.ProposalType
	// Ledger only
	derivationPath string
)

var SignPrivateKeyCmd = &cobra.Command{
	Use:   "sign-pk",
	Short: "Sign a proposal with a private key",
	Long:  `Configure a private key in a .env file (using the PRIVATE_KEY var) and sign a proposal with it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get private key
		pk, err := LoadPrivateKey()
		if err != nil {
			return err
		}

		// Load proposal
		p, err := proposal.LoadProposal(proposalType, proposalPath)

		if err != nil {
			return err
		}

		err = proposal.SignPlainKey(pk, p)
		if err != nil {
			return err
		}

		// Write proposal to file
		err = WriteProposalToFile(p, proposalPath)
		if err != nil {
			return err
		}

		return nil
	},
}

func WriteProposalToFile(proposal interface{}, filePath string) error {
	proposalBytes, err := json.Marshal(proposal)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, proposalBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

var SignLedgerCmd = &cobra.Command{
	Use:   "sign-ledger",
	Short: "Sign a proposal with a ledger",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load proposal
		proposal, err := proposal.LoadProposal(proposalType, proposalPath)
		if err != nil {
			return err
		}

		// Parse the derivation path
		path, err := accounts.ParseDerivationPath(derivationPath)
		if err != nil {
			log.Fatalf("Failed to parse derivation path: %v", err)
		}

		err = SignLedger(path, proposal)
		if err != nil {
			return err
		}
		//
		//// Write proposal to file
		//err = WriteProposalToFile(proposal, proposalPath)
		//if err != nil {
		//	return err
		//}
		//
		return nil
	},
}

func LoadPrivateKey() (*ecdsa.PrivateKey, error) {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	// Load PrivateKey
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		return nil, errors.New("PRIVATE_KEY not found in .env file")
	}

	// Convert to ecdsa
	ecdsa, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}
	return ecdsa, nil
}

// Just run this locally to sign from the ledger.
func SignLedger(derivationPath []uint32, proposal proposal.Proposal) error {
	// Validate proposal
	err := proposal.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate proposal: %w", err)
	}

	// Load ledger
	ledgerhub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return fmt.Errorf("failed to open ledger hub: %w", err)
	}

	// Get the first wallet
	wallets := ledgerhub.Wallets()
	if len(wallets) == 0 {
		return errors.New("no wallets found")
	}
	wallet := wallets[0]

	fmt.Printf("Found %d wallets, using first one\n",
		len(wallets))

	// Open the ledger.
	err = wallet.Open("")
	if err != nil {
		return fmt.Errorf("failed to open wallet: %w", err)
	}

	fmt.Printf("Opened wallet, have accounts %v\n", wallet.Accounts())

	// Load account.
	// BIP44 derivation path used in ledger.
	// Could pass this in as an argument as well.
	account, err := wallet.Derive(derivationPath, true)
	if err != nil {
		return fmt.Errorf("failed to derive account: %w derivation path %v"+
			" is your ledger ethereum app open?", err, derivationPath)
	}
	fmt.Println("Found account: ", account.Address.String())

	// Get the signing hash

	// Create executor
	executor, err := proposal.ToExecutor(false)
	if err != nil {
		return err
	}

	payload, err := executor.SigningHash()
	if err != nil {
		return err
	}

	// Sign the payload
	sig, err := wallet.SignData(account, accounts.MimetypeTypedData, payload.Bytes())
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

func main() {
	rootCmd := cobra.Command{}
	rootCmd.PersistentFlags().StringVar(&proposalPath, "proposal", "p", "Absolute file path containing the proposal to be submitted")
	var proposalTypeStr string
	rootCmd.PersistentFlags().StringVar(&proposalTypeStr, "proposalType", string(proposal.MCMS), "The type of proposal being ingested")
	proposalType = proposal.StringToProposalType[proposalTypeStr]

	SignLedgerCmd.PersistentFlags().StringVar(&derivationPath, "derivationPath", "m/44'/60'/0'/0/0", "The type of proposal being ingested")
	rootCmd.AddCommand(SignLedgerCmd)
	rootCmd.AddCommand(SignPrivateKeyCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
