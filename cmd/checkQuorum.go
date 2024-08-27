package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var checkQuorumCmd = &cobra.Command{
	Use:   "check-quorum",
	Short: "Determines whether the provided signatures meet the quorum to set the root",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := BuildExecutableProposal(proposal)
		if err != nil {
			log.Fatal(err)
		}

		clientBackend, err := BuildContractDeployBackend()
		if err != nil {
			log.Fatal(err)
		}

		executor, err := executableProposal.ToExecutor(clientBackend)
		if err != nil {
			log.Fatal(err)
		}

		err = executor.ValidateSignatures()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkQuorumCmd)
}
