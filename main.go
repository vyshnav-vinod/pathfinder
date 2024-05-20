package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		success(os.Stdout, path, c)
	}
	if path == "" {
		flaggy.ShowHelpAndExit("Please provide arguments")
	}

	pathfinder(os.Stdout, c, ignoreDir, path)
}

func pathfinder(w io.Writer, c *Cache, ignoreDir bool, path string) {
	cleanedPath, err := filepath.Abs(path)
	HandleError(err)
	if checkDirExists(cleanedPath) {
		// Abs might join path with cwd, so this will
		// also check if the directory is in the cwd
		if ignoreDir {
			cwd, err := os.Getwd()
			HandleError(err)
			if !(filepath.Dir(cleanedPath) == cwd) {
				success(w, cleanedPath, c)
			}
		} else {
			success(w, cleanedPath, c)
		}
	}

	// Assume that the path is not an actual path but a search query by the user and it might exist
	if cacheEntry, ok := c.GetCacheEntry(filepath.Base(path)); ok {
		success(w, cacheEntry.Path, c)
	}
	var found = false
	if !ignoreDir {
		traverseAndMatchDir(w, ".", path, &found, c) // Walk inside working directory
	}
	cwd, err := os.Getwd()
	HandleError(err)
	traverseAndMatchDir(w, filepath.Dir(cwd), path, &found, c) // Walk from one directory above
	usrHome, err := os.UserHomeDir()
	HandleError(err)
	traverseAndMatchDir(w, usrHome, path, &found, c) // Walk from $HOME

	if !found {
		// pathfinder failed to find the directory and prints
		// the path (user input) to stdout for the bash script
		// to capture and return as an error msg
		// fmt.Println(path)
		fmt.Fprint(w, path)
		os.Exit(EXIT_FOLDERNOTFOUND)
	}

}

func traverseAndMatchDir(w io.Writer, dirName string, searchDir string, found *bool, c *Cache) {
	if !strings.HasPrefix(filepath.Base(dirName), ".") {
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
					*found = true
					success(w, path, c)
				} else {
					traverseAndMatchDir(w, path, searchDir, found, c)
				}
			}
		}
	}
}

func success(w io.Writer, path string, c *Cache) {
	// Prints to stdout for bash script to capture
	// fmt.Println(path)
	fmt.Fprint(w, path)
	c.SetPreviousDir()
	c.SetCacheEntry(path)
	os.Exit(EXIT_SUCCESS)
}

func checkDirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
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
