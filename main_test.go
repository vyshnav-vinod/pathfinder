package main

import (
	"os"
	"testing"
)

var tmpdir string

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func teardown() {
	HandleError(os.Unsetenv("PF_TMP_TEST"))
}

// While testing, a tmp folder will be made.
// It will contain some sub folders as well.
// The plan for the test of the main function
func TestMainFunc(t *testing.T) {
}
