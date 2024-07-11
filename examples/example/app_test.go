package main

import (
	"os"
	"testing"
)

func TestExampleAppExecution(t *testing.T) {
	originalArgs := os.Args
	// Ensure os.Args is restored after this test
	defer func() { os.Args = originalArgs }()

	// No arguments
	os.Args = []string{""}
	RunMain()
	t.Log("Example ran successfully")
}
