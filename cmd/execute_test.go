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
	rootCmd.SetArgs([]string{"execute", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1", "--pk", "0x123"})
	rootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", rpc)
	assert.Equal(t, "./proposal.json", proposalPath)
	assert.Equal(t, "1", chainSelector)
	assert.Equal(t, "0x123", pk)

	expectedDescription := "help for execute"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
