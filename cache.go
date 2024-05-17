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

// The json will be of the type map[stringCacheSchema
type CacheSchema struct {
	Path      string `json:"path"`
	Frequency int    `json:"frequency"`
}

type Cache struct {
	file     string
	maxCap   int
	contents map[string]CacheSchema
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

func (c *Cache) LoadCache() {
	// Load the cache contents to memory
	f, err := os.ReadFile(c.file)
	HandleError(err)
	HandleError(json.Unmarshal(f, &c.contents))
}

func (c *Cache) GetPreviousDir() string {
	if _, ok := c.contents[PREV_DIR_ENTRY]; ok {
		return c.contents[PREV_DIR_ENTRY].Path
	} else {
		c.SetPreviousDir()
	return c.contents[PREV_DIR_ENTRY].Path
}
}

func (c *Cache) SetPreviousDir() {
	cwd, err := os.Getwd()
	HandleError(err)
	c.contents[PREV_DIR_ENTRY] = CacheSchema{
		Path:      cwd,
		Frequency: -1,
	}
	writeContent, err := json.MarshalIndent(c.contents, "", " ")
	HandleError(err)
	HandleError(os.WriteFile(c.file, writeContent, 077))
}

func (c *Cache) GetCacheEntry(entry string, path string) (cacheEntry CacheSchema, ok bool) {
	if cache, found := c.contents[entry]; found {
		return cache, true
	}
	return CacheSchema{}, false
}

// func (c *Cache) SetCacheEntry(entry string) {
// Shoudl check if entry is in cache ,if yes just update
// else pop from cache and add new entry
// }

// func (c *Cache) popCache() {
// 	// Removes according to LFU
// }

// func (c *Cache) clearCache() {
// 	// Optional
// }
