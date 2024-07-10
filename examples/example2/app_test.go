package main

import (
	"os/exec"
	"testing"
)

// Build and run the application to test it.
// Any issues in the example should return non-zero exit code.
func TestAppExecution(t *testing.T) {
	cmdBuild := exec.Command("go", "build", "-o", "./example")
	err := cmdBuild.Run()
	if err != nil {
		t.Fatalf("Building failed: %v", err)
	}
	cmdRun := exec.Command("./example")
	err = cmdRun.Run()
	if err != nil {
		t.Fatalf("Execution failed: %v", err)
	}
	// If the app exits with a non-zero exit code, cmd.Run() will return an error.
	// If we reach this point, the app has executed successfully with a zero exit code.
	t.Log("Example ran successfully")
}
