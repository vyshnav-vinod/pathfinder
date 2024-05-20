package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const testEnv = "PF_TMP_TEST"

func TestMainFunc(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.Setenv(testEnv, tmpDir)
	if err != nil {
		t.Fatalf("Error setting test environment key : %v", err)
	}

	defer t.Cleanup(func() {
		HandleError(os.Unsetenv(testEnv))
	})

	// creating subfolders inside tmpDir
	var subTmpDirs [5]string
	for i := range subTmpDirs {
		fname := filepath.Join(tmpDir, fmt.Sprintf("subdir%d", i))
		err := os.Mkdir(fname, 0777)
		if err != nil {
			t.Fatalf("Error while creating temporary subfolders : %v", err)
		}
		subTmpDirs[i] = fname
	}

	/* Test Cases
		1) Valid dir change 1
		2) Valid dir change 2
		3) Dir not found
	*/
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
