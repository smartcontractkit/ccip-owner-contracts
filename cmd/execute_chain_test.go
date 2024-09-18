package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteChainCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"execute-chain", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1"})
	rootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", rpc)
	assert.Equal(t, "./proposal.json", proposalPath)
	assert.Equal(t, uint64(1), chainSelector)

	expectedDescription := "help for execute"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
