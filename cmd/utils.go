package cmd

import (
	"encoding/json"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

func BuildExecutableProposal(path string) (*mcms.Proposal, error) {
	var proposal = mcms.Proposal{}
	raw, err := readFromFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &proposal)
	return &proposal, err
}

func readFromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
