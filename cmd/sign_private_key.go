package cmd

import (
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/spf13/cobra"
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
		p, err := LoadProposal(proposalType, proposalPath)
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
