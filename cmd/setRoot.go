package cmd

import (
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	cs "github.com/smartcontractkit/chain-selectors"
)

var setMerkleCmd = &cobra.Command{
	Use:   "set-root",
	Short: "Sets the Merkle Root on the MCM Contract",
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

		auth, err := bind.NewKeyedTransactorWithChainID(ecdsa, bigIntChainID)
		if err != nil {
			log.Fatal(err)
		}

		transaction, err := executor.SetRootOnChain(auth, chainSelector)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Transaction sent: %s", transaction.Hash().Hex())
	},
}

func init() {
	rootCmd.AddCommand(setMerkleCmd)
	setMerkleCmd.PersistentFlags().StringVar(&pk, "pk", "0", "Private key used to send the transaction")
}
