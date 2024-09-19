package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetRootCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs([]string{"set-root", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1"})
	RootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", Rpc)
	assert.Equal(t, "./proposal.json", ProposalPath)
	assert.Equal(t, uint64(1), ChainSelector)

	expectedDescription := "no such file or directory"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
