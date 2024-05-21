package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/integrii/flaggy"
)

const (
	EXIT_SUCCESS        = 0
	EXIT_FOLDERNOTFOUND = 1
	EXIT_CACHECLEANED   = 4
	EXIT_ERR            = -1
)

// TODO: Add a version and build number
func main() {

	// Flags
	var (
		ignoreDir   = false
		previousDir = false
		cleanCache  = false
		path        string
	)

	flaggy.Bool(&ignoreDir, "i", "ignore", "Ignore searching the current directory")
	flaggy.Bool(&previousDir, "b", "back", "Change directory to the previous directory")
	flaggy.Bool(&cleanCache, "", "clean", "Clean the cache")
	flaggy.AddPositionalValue(&path, "Directory", 1, false, "The name/path of the directory")
	flaggy.Parse()

	c := InitCache()

	if cleanCache {
		c.cleanCache()
		os.Exit(EXIT_CACHECLEANED)
	}
	if previousDir {
		path, _ := filepath.Abs(c.GetPreviousDir())
		os.Exit(success(os.Stdout, path, c))
	}
	if path == "" {
		flaggy.ShowHelpAndExit("Please provide arguments")
	}

	os.Exit(pathfinder(os.Stdout, c, ignoreDir, path))
}

func pathfinder(w io.Writer, c *Cache, ignoreDir bool, path string) int {

	absPath, err := filepath.Abs(path)
	HandleError(err)
	cwd, err := os.Getwd()
	HandleError(err)

	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		// To support ~, .. , etc
		if !ignoreDir {
			return success(w, absPath, c)
		} else {
			if !strings.Contains(absPath, cwd) { // Check if it is in current directory
				return success(w, absPath, c)
			}
		}
	}

	// Assume that the path is not an actual path but a search query by the user and it might exist
	if cacheEntry, ok := c.GetCacheEntry(filepath.Base(path)); ok {
		if !ignoreDir {
			return success(w, cacheEntry.Path, c)
		} else {
			if !strings.Contains(cacheEntry.Path, cwd) { // Check if it is in current directory
				return success(w, cacheEntry.Path, c)
			}
		}
	}

	var pathReturned string
	var dirsAlreadyWalked []string // to ignore walking through already walked directories

	// TODO: Goroutines or a new algorithm
	if !ignoreDir {
		if traverseAndMatchDir(w, cwd, path, &pathReturned, dirsAlreadyWalked, c) {
			// Walk inside working directory
			return success(w, pathReturned, c)
		}
	}

	dirsAlreadyWalked = append(dirsAlreadyWalked, cwd)
	if traverseAndMatchDir(w, filepath.Dir(cwd), path, &pathReturned, dirsAlreadyWalked, c) {
		// Walk from one directory above
		return success(w, pathReturned, c)
	}

	usrHome, err := os.UserHomeDir()
	HandleError(err)
	dirsAlreadyWalked = append(dirsAlreadyWalked, filepath.Dir(cwd))
	if traverseAndMatchDir(w, usrHome, path, &pathReturned, dirsAlreadyWalked, c) {
		// Walk from $HOME
		return success(w, pathReturned, c)
	}

	// pathfinder failed to find the directory and prints
	// the path (user input) to stdout for the bash script
	// to capture and return as an error msg
	fmt.Fprint(w, path)
	return EXIT_FOLDERNOTFOUND
}

func traverseAndMatchDir(w io.Writer, dirName string, searchDir string, pathReturned *string, dirsAlreadyWalked []string, c *Cache) bool {
	if !strings.HasPrefix(filepath.Base(dirName), ".") && !slices.Contains(dirsAlreadyWalked, dirName) {
		file, err := os.Open(dirName)
		HandleError(err)
		defer file.Close()
		dirEntries, err := file.Readdirnames(0)
		HandleError(err)
		for _, n := range dirEntries {
			path, err := filepath.Abs(filepath.Join(dirName, n))
			HandleError(err)
			f, err := os.Stat(path)
			if os.IsNotExist(err) {
				continue
			}
			if f.IsDir() {
				if f.Name() == searchDir {
					*pathReturned = path
					return true
				} else {
					if traverseAndMatchDir(w, path, searchDir, pathReturned, dirsAlreadyWalked, c) {
						return true
					}
				}
			}
		}
	}
	return false
}

func success(w io.Writer, path string, c *Cache) int {
	// Prints to stdout for bash script to capture
	fmt.Fprint(w, path)
	c.SetPreviousDir()
	c.SetCacheEntry(path)
	return EXIT_SUCCESS
}

func HandleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %s\n", err)
		os.Exit(EXIT_ERR)
	}
}

func InitCache() *Cache {
	var cacheFile string
	if _, ok := os.LookupEnv("PF_TMP_TEST"); !ok {
		// This is done for testing. See main_test.go for more info
		cf, err := os.UserCacheDir()
		HandleError(err)
		cacheFile = filepath.Join(cf, "pathfinder", "cache.json")
	} else {
		// This is done for testing. See main_test.go for more info
		cacheFile = filepath.Join(os.Getenv("PF_TMP_TEST"), "cache.json")
	}
	c := &Cache{
		file:   cacheFile,
		maxCap: 10,
	}
	c.CheckCache()
	c.LoadCache()
	return c
}
