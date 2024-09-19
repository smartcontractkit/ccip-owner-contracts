package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cliRootCommand(t *testing.T) {

	actual := new(bytes.Buffer)
	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs([]string{})
	RootCmd.Execute()

	expectedDescription := "Tools for on-chain interactions with the MCMS"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}
