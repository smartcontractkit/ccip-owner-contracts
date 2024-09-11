package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetRootCommand(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"set-root", "--rpc", "http://localhost:8545", "--proposal", "./proposal.json", "--selector", "1", "--pk", "0x123"})
	rootCmd.Execute()

	assert.Equal(t, "http://localhost:8545", rpc)
	assert.Equal(t, "./proposal.json", proposalPath)
	assert.Equal(t, "1", chainSelector)
	assert.Equal(t, "0x123", pk)

	expectedDescription := "no such file or directory"
	assert.Containsf(t, actual.String(), expectedDescription, "expected description to contain '%s'", expectedDescription)
}

func Test_SetRootCommandWithFile(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)

	temp, err := os.CreateTemp("", "testName")
	defer os.Remove(temp.Name())
	if err != nil {
		t.Fatal(err)
	}

	_, err = temp.Write([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}

	// update this test when we have better proposal validation in place
	rootCmd.SetArgs([]string{"set-root", "--rpc", "http://localhost:8545", "--proposal", temp.Name(), "--selector", "14767482510784806043", "--pk", "0x123"})
	shouldPanic(t, rootCmd.Execute)

	assert.Equal(t, "14767482510784806043", chainSelector)
}

func shouldPanic(t *testing.T, f func() error) {
	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}
