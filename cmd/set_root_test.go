package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetRootCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"set-root", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1"})
	rootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", rpc)
	assert.Equal(t, "./proposal.json", proposalPath)
	assert.Equal(t, uint64(1), chainSelector)

	expectedDescription := "no such file or directory"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
