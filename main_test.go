package main

import (
	"os"
	"testing"
)

const testEnv = "PF_TMP_TEST"

func Test_pathfinder(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.Setenv(testEnv, tmpDir)
	if err != nil {
		t.Fatalf("Error setting test environment key : %v", err)
	}

	defer t.Cleanup(func() {
		HandleError(os.Unsetenv(testEnv))
	})

	// TODO: Tests
}
