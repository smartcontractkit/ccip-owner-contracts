package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"execute-operation", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1", "--index", "3"})
	rootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", rpc)
	assert.Equal(t, "./proposal.json", proposalPath)
	assert.Equal(t, uint64(1), chainSelector)
	assert.Equal(t, uint64(3), index)

	expectedDescription := "help for execute"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
