package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/integrii/flaggy"
)

const (
	EXIT_SUCCESS        = 0
	EXIT_FOLDERNOTFOUND = 1
	EXIT_ERR            = -1
)

func main() {

	exitCode := EXIT_SUCCESS
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

	startTime := time.Now()
	cleanedPath, err := filepath.Abs(path)
	handleError(err)
	if checkDirExists(cleanedPath) {
		// Abs might join path with cwd, so this will
		// also check if the directory is in the cwd
		fmt.Println(cleanedPath)
	} else {
		// Assume that the path is not an actual path but a search query by the user and it might exist

		// TODO: Check the cache first and then proceed. (Make a seperate cache function to check for entry)

		var returnedPath = ""
		traverseAndMatchDir(".", path, &returnedPath)
		cwd, err := os.Getwd()
		handleError(err)
		traverseAndMatchDir(filepath.Dir(cwd), path, &returnedPath)
		usrHome, err := os.UserHomeDir()
		handleError(err)
		traverseAndMatchDir(usrHome, path, &returnedPath)

		if len(returnedPath) == 0 {
			fmt.Println(path)
			os.Exit(EXIT_FOLDERNOTFOUND)
		} else {
			fmt.Println(returnedPath)
			os.Exit(EXIT_SUCCESS)
		}
	}
	fmt.Printf("it took %v \n", time.Since(startTime)) // defer
	// Max time : 1.6s (Search not found)
}

func traverseAndMatchDir(dirName string, searchDir string, pathReturn *string) bool {
	file, err := os.Open(dirName)
	handleError(err)
	defer file.Close()
	dirEntries, err := file.Readdirnames(0)
	handleError(err)
	for _, n := range dirEntries {
		path, err := filepath.Abs(filepath.Join(dirName, n))
		handleError(err)
		f, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}
		if f.IsDir() {
			if f.Name() == searchDir {
				*pathReturn = path
				return true
			} else {
				if traverseAndMatchDir(path, searchDir, pathReturn) {
					return true
				}
			}
		}
	}
	return false
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
		os.Exit(EXIT_ERR)
	}
}
