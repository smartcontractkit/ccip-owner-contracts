package cmd

import (
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	cs "github.com/smartcontractkit/chain-selectors"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Performs an operation execution on the MCMS. Root must be set first.",
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

		uintChainSelector, err := strconv.ParseUint(chainSelector, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		// retrieve current chainID (can't use chain selector, unless?)
		chainID, err := cs.ChainIdFromSelector(uintChainSelector)
		if err != nil {
			log.Fatalf("Failed to retrieve chain ID from chain selector %v", err)
		}

		ecdsa, err := crypto.HexToECDSA(pk)
		if err != nil {
			log.Fatal(err)
		}

		bigIntChainID := big.NewInt(int64(chainID))

		for i, op := range executableProposal.Transactions {
			if op.ChainIdentifier == chainSelector {
				auth, err := bind.NewKeyedTransactorWithChainID(ecdsa, bigIntChainID)
				if err != nil {
					log.Fatal(err)
				}

				transaction, err := executor.ExecuteOnChain(auth, i)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Update pending: 0x%x\n", transaction.Hash())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.PersistentFlags().StringVar(&pk, "pk", "0", "Private key used to send the transaction")
}
