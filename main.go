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

	// Flags
	var ignoreDir = false
	var previousDir = false
	var path string

	flaggy.Bool(&ignoreDir, "i", "ignore", "Ignore searching the current directory")
	flaggy.Bool(&previousDir, "b", "back", "Change directory to the previous directory")
	flaggy.AddPositionalValue(&path, "Directory", 1, false, "The name/path of the directory")
	flaggy.Parse()

	startTime := time.Now()
	c := InitCache()
	if previousDir {
		path, _ := filepath.Abs(c.GetPreviousDir())
		success(path, c)
	}
	cleanedPath, err := filepath.Abs(path)
	HandleError(err)
	if checkDirExists(cleanedPath) {
		// Abs might join path with cwd, so this will
		// also check if the directory is in the cwd
		if ignoreDir {
			cwd, err := os.Getwd()
			HandleError(err)
			if !(filepath.Dir(cleanedPath) == cwd) {
				success(cleanedPath, c)
			}
		} else {
			success(cleanedPath, c)
		}
	}

	// Assume that the path is not an actual path but a search query by the user and it might exist

	var returnedPath = ""
	if !ignoreDir {
		traverseAndMatchDir(".", path, &returnedPath)
	}
	cwd, err := os.Getwd()
	HandleError(err)
	traverseAndMatchDir(filepath.Dir(cwd), path, &returnedPath)
	usrHome, err := os.UserHomeDir()
	HandleError(err)
	traverseAndMatchDir(usrHome, path, &returnedPath)

	if len(returnedPath) == 0 {
		fmt.Println(path)
		os.Exit(EXIT_FOLDERNOTFOUND)
	}
	fmt.Printf("it took %v \n", time.Since(startTime)) // make it deferred
	// Max time : 1.6s (Search not found)
}

func traverseAndMatchDir(dirName string, searchDir string, pathReturn *string) bool {
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
				*pathReturn = path
				fmt.Println(*pathReturn)
				os.Exit(EXIT_SUCCESS)
			} else {
				traverseAndMatchDir(path, searchDir, pathReturn)
			}
		}
	}
	return false
}

func success(path string, c *Cache){
	fmt.Println(path)
	c.WritePreviousDir()
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
	cf, err := os.UserCacheDir()
	HandleError(err)
	cacheFile := filepath.Join(cf, "pathfinder", "cache.json")
	c := &Cache{
		file:   cacheFile,
		maxCap: 10,
	}
	c.CheckCache()
	c.LoadCache()
	return c
}
