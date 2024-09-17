package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rpc string
var proposalPath string
var chainSelector string
var pk string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcms",
	Short: "Tools for on-chain interactions with the MCMS ",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.AddCommand(CheckQuorumCmd)
	rootCmd.AddCommand(ExecuteOperationCmd)
	rootCmd.AddCommand(SetMerkleCmd)

	rootCmd.PersistentFlags().StringVar(&rpc, "rpc", "", "rpc to be used in the proposal")
	rootCmd.PersistentFlags().StringVar(&proposalPath, "proposal", "p", "Absolute file path containing the proposal to be submitted")
	rootCmd.PersistentFlags().StringVar(&chainSelector, "selector", "-1", "Chain selector for the command to connect to")
}
