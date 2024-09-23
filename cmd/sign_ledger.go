package cmd

/*
import (
	"github.com/smartcontractkit/ccip-owner-contracts/tools/signing"
	"github.com/spf13/cobra"
)

var derivationPath []uint

var SignLedgerCmd = &cobra.Command{
	Use:   "sign-pk",
	Short: "Sign a proposal with a private key",
	Long:  `Configure a private key in a .env file (using the PRIVATE_KEY var) and sign a proposal with it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load proposal
		proposal, err := LoadProposal(proposalType, proposalPath)
		if err != nil {
			return err
		}

		// convert derivation path to []uint32
		derivationPathUint32 := make([]uint32, len(derivationPath))
		for i, v := range derivationPath {
			derivationPathUint32[i] = uint32(v)
		}

		err = signing.SignLedger(derivationPathUint32, proposal)
		if err != nil {
			return err
		}

		// Write proposal to file
		err = WriteProposalToFile(proposal, proposalPath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	SignLedgerCmd.Flags().UintSliceVar(&derivationPath, "derivation-path", []uint{44, 60, 0, 0, 0}, "Derivation path for the Ledger")
}
*/
