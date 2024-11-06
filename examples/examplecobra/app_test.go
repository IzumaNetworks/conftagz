package main

import (
	"os"
	"testing"
)

// Build and run the application to test it.
// Any issues in the example should return non-zero exit code.
func TestExampleCobraExecution(t *testing.T) {
	originalArgs := os.Args
	// Ensure os.Args is restored after this test
	defer func() { os.Args = originalArgs }()
	// No arguments - original arguments can contain some test configs,
	// which messes the test
	os.Args = []string{"progname", "othercmd", "--anotherfield", "ApplesAreGreat"}
	RunMain()
	t.Log("Examplecobra ran successfully")
}
