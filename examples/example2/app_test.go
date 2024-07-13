package main

import (
	"os"
	"testing"
)

// Build and run the application to test it.
// Any issues in the example should return non-zero exit code.
func TestExample2Execution(t *testing.T) {
	originalArgs := os.Args
	// Ensure os.Args is restored after this test
	defer func() { os.Args = originalArgs }()
	// No arguments - original arguments can contain some test configs,
	// which messes the test
	os.Args = []string{""}

	RunMain()
	t.Log("Example2 ran successfully")
}
