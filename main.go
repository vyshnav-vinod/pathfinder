package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/integrii/flaggy"
)

const (
	EXIT_SUCCESS        = 0
	EXIT_FOLDERNOTFOUND = 1
	EXIT_CACHECLEANED   = 4
	EXIT_INFO           = 5
	EXIT_ERR            = -1
)

var (
	VERSION     string
	BUILD_NUM   string
	NAME        = "pf"
	DESCRIPTION = "A command line tool to move between directories easily and fast."
)

var (
	// go run -ldflags "-X main.DEV=true" .
	DEV       = "false" // Ignore, Build tag used for dev testing and some benchmarks
	timeStart time.Time
)

type pathInfo struct {
	// Stores the user input
	userInput string
	// Restrict is used to denote whether user has specified
	// the parent directory of the directory to search as well.
	// Eg: User has specified "parent/searchfolder"
	// We must only find "searchfolder" which is a child of the
	// directory named "parent". Any other folders named
	// "searchfolder" will be rejected if their parent directory
	// is not "parent"
	restrict bool
	// Full path (if found)
	path string
}

func main() {

	// Ignore, This is a build flag for testing purposes
	if DEV == "true" {
		timeStart = time.Now()
	}

	defer func() {
		if DEV == "true" {
			fmt.Printf("It took %v\n", time.Since(timeStart))
		}
	}()

	// Flags
	var (
		ignoreDir   = false
		previousDir = false
		cleanCache  = false
		info        = false
		path        string
	)

	flaggy.Bool(&ignoreDir, "i", "ignore", "Ignore searching the current directory")
	flaggy.Bool(&previousDir, "b", "back", "Change directory to the previous directory")
	flaggy.Bool(&cleanCache, "", "clean", "Clean the cache")
	flaggy.Bool(&info, "", "info", "Display version and build number")
	flaggy.AddPositionalValue(&path, "Directory", 1, false, "The name/path of the directory")
	flaggy.SetName(NAME)
	flaggy.SetDescription(DESCRIPTION)
	flaggy.DefaultParser.DisableShowVersionWithVersion()
	flaggy.Parse()

	if info {
		fmt.Fprintf(os.Stdout, "Version:%s  Build:%s", VERSION, BUILD_NUM)
		os.Exit(EXIT_INFO)
	}

	c := InitCache()

	if cleanCache {
		c.cleanCache()
		os.Exit(EXIT_CACHECLEANED)
	}
	if previousDir {
		path, _ := filepath.Abs(c.GetPreviousDir())
		os.Exit(success(os.Stdout, pathInfo{path: path, restrict: false, userInput: filepath.Base(path)}, c))
	}
	if path == "" {
		flaggy.ShowHelpAndExit("Please provide arguments")
	}

	// Ignore, This is a build flag for testing purposes
	if DEV == "false" {
		os.Exit(pathfinder(os.Stdout, c, ignoreDir, path))
	} else {
		fmt.Printf("\npathfinder returned exit code %d\n", pathfinder(os.Stdout, c, ignoreDir, path))
	}
}

func pathfinder(w io.Writer, c *Cache, ignoreDir bool, searchPath string) int {

	absPath, err := filepath.Abs(searchPath)
	if err != nil {
		HandleError(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		HandleError(err)
	}

	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		// To support ~, .. , etc
		if !ignoreDir {
			return success(w, pathInfo{path: absPath, restrict: false, userInput: filepath.Base(absPath)}, c)
		} else {
			if !strings.Contains(absPath, cwd) { // Check if it is in current directory
				return success(w, pathInfo{path: absPath, restrict: false, userInput: filepath.Base(absPath)}, c)
			}
		}
	}

	pathInfo := pathInfo{userInput: searchPath}
	pathInfo.restrict = !(len(strings.Split(searchPath, "/")) == 1)

	// Assume that the path is not an actual path but a search query by the user and it might exist
	if cacheEntry, ok := c.GetCacheEntry(pathInfo); ok {
		if !ignoreDir {
			pathInfo.path = cacheEntry.Path
			return success(w, pathInfo, c)
		} else {
			if !strings.Contains(cacheEntry.Path, cwd) { // Check if it is in current directory
				pathInfo.path = cacheEntry.Path
				return success(w, pathInfo, c)
			}
		}
	}

	dirsAlreadyWalked := make(map[string]struct{}) // to ignore walking through already walked directories

	// TODO: Goroutines or a new algorithm
	if !ignoreDir {
		if traverseAndMatchDir(w, cwd, pathInfo, &pathInfo.path, dirsAlreadyWalked, c) {
			// Walk inside working directory
			return success(w, pathInfo, c)
		}
	}

	dirsAlreadyWalked[cwd] = struct{}{}
	if traverseAndMatchDir(w, filepath.Dir(cwd), pathInfo, &pathInfo.path, dirsAlreadyWalked, c) {
		// Walk from one directory above
		return success(w, pathInfo, c)
	}

	usrHome, err := os.UserHomeDir()
	if err != nil {
		HandleError(err)
	}
	dirsAlreadyWalked[filepath.Dir(cwd)] = struct{}{}
	if traverseAndMatchDir(w, usrHome, pathInfo, &pathInfo.path, dirsAlreadyWalked, c) {
		// Walk from $HOME
		return success(w, pathInfo, c)
	}

	// pathfinder failed to find the directory and prints
	// the path (user input) to stdout for the bash script
	// to capture and return as an error msg
	fmt.Fprint(w, searchPath)
	return EXIT_FOLDERNOTFOUND
}

func traverseAndMatchDir(w io.Writer, dirName string, searchDir pathInfo, pathReturned *string, dirsAlreadyWalked map[string]struct{}, c *Cache) bool {
	if strings.HasPrefix(filepath.Base(dirName), ".") {
		return false
	}
	if _, ok := dirsAlreadyWalked[dirName]; ok {
		return false
	}
	file, err := os.Open(dirName)
	if err != nil {
		HandleError(err)
	}
	defer file.Close()
	dirEntries, err := file.Readdirnames(0)
	if err != nil {
		HandleError(err)
	}
	for _, n := range dirEntries {
		path, err := filepath.Abs(filepath.Join(dirName, n))
		if err != nil {
			HandleError(err)
		}
		f, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}
		if f.IsDir() {
			if f.Name() == filepath.Base(searchDir.userInput) {
				if !searchDir.restrict {
					*pathReturned = path
					return true
				} else {
					if strings.Contains(path, searchDir.userInput) {
						*pathReturned = path
						return true
					} else {
						if traverseAndMatchDir(w, path, searchDir, pathReturned, dirsAlreadyWalked, c) {
							return true
						}
					}
				}
			} else {
				if traverseAndMatchDir(w, path, searchDir, pathReturned, dirsAlreadyWalked, c) {
					return true
				}
			}
		}
	}

	return false
}

func success(w io.Writer, path pathInfo, c *Cache) int {
	// Prints to stdout for bash script to capture
	fmt.Fprint(w, path.path)
	c.SetPreviousDir()
	fmt.Println("\nentry : ", path)
	c.SetCacheEntry(path)
	return EXIT_SUCCESS
}

func HandleError(err error) {
	fmt.Fprintf(os.Stderr, "[Error] %s\n", err)
	os.Exit(EXIT_ERR)
}

func InitCache() *Cache {
	var cacheFile string
	if _, ok := os.LookupEnv("PF_TMP_TEST"); !ok {
		// This is done for testing. See main_test.go for more info
		cf, err := os.UserCacheDir()
		if err != nil {
			HandleError(err)
		}
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
