/*
	cache is stored in USERCACHEDIR/pathfinder/

	Structure of the cache
	{
		"filename" : {
			"path" : ...,
			"frequency" : ...
		}
	}
	The first entry will always be the previous directory redirected by pathfinder
	And the value of "filename" will be PFpreviousDir
*/

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Cache struct {
	file   string
	maxCap int
}

// The json will be of the type map[stringCacheSchema
type CacheSchema struct {
	Path      string `json:"path"`
	Frequency int    `json:"frequency"`
}

const PREV_DIR_ENTRY = "PFpreviousDir"

func (c *Cache) CheckCache() {
	// Check if the cache file is valid
	// If there is no cache file, make a new cache file

	// Cache.file should be cache.json in the user's cache dir
	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		HandleError(os.MkdirAll(filepath.Dir(c.file), 0777))
		_, err = os.Create(c.file)
		HandleError(err)
	}
}

func (c *Cache) GetPreviousDir() string {
	f, err := os.ReadFile(c.file)
	HandleError(err)
	var cacheMap map[string]CacheSchema
	HandleError(json.Unmarshal(f, &cacheMap))
	return cacheMap[PREV_DIR_ENTRY].Path
}
