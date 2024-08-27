package cmd

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	e "github.com/smartcontractkit/ccip-owner-contracts/tools/executable"
)

func BuildExecutableProposal(path string) error {
	raw, err := readFromFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &executableProposal)
	return err
}

func readFromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func BuildContractDeployBackend() (map[string]e.ContractDeployBackend, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	return map[string]e.ContractDeployBackend{
		chainSelector: client,
	}, nil
}
