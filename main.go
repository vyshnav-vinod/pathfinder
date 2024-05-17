package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/integrii/flaggy"
)

func main() {

	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	// Flags
	var previousDir = false
	var ignoreDir = false
	var path string

	flaggy.Bool(&ignoreDir, "i", "ignore", "Ignore searching the current directory")
	flaggy.Bool(&previousDir, "b", "back", "Change directory to the previous directory")
	flaggy.AddPositionalValue(&path, "Directory", 1, true, "The name/path of the directory")
	flaggy.Parse()

	cleanedPath, err := filepath.Abs(path)
	handleError(err)
	if checkDirExists(cleanedPath) {
		fmt.Println(cleanedPath)
	} else {
		exitCode = 33
	}
}

func checkDirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %s\n", err)
		os.Exit(-1)
	}
}
