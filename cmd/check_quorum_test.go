package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckQuorumCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs([]string{"check-quorum", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1"})
	RootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", Rpc)
	assert.Equal(t, "./proposal.json", ProposalPath)
	assert.Equal(t, uint64(1), ChainSelector)

	expectedDescription := "help for check-quorum"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
