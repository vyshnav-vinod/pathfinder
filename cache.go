/*
	cache is stored in USERCACHEDIR/pathfinder/

	Structure of the cache
	{
		"filename" : {
			"path" : ...,
			"frequency" : ...,
			"lasthit": ...
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
	"time"
)

// The json will be of the type map[string]CacheSchema
type CacheSchema struct {
	Path      string    `json:"path"`
	Frequency int       `json:"frequency"`
	LastHit   time.Time `json:"lasthit"`
}

type Cache struct {
	file     string
	maxCap   int // Only the capacity of the cache items, Does not consider additionals like previous dir entry
	contents map[string]CacheSchema
}

const PREV_DIR_ENTRY = "PFpreviousDir"

func (c *Cache) CheckCache() {
	// Check if the cache file exist
	// If there is no cache file, make a new cache file

	// Cache.file should be cache.json in the user's cache dir
	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		HandleError(os.MkdirAll(filepath.Dir(c.file), 0777))
		_, err = os.Create(c.file)
		HandleError(err)
	}
	c.validateCache()
}

func (c *Cache) validateCache() {
	// Check if cache file is empty and if it is
	// write a default cache data into it
	// This check is to prevent "assignment to entry in nil map panic"
	f, err := os.ReadFile(c.file)
	HandleError(err)
	if len(f) == 0 {
		tmpMap := make(map[string]CacheSchema)
		usrHome, err := os.UserHomeDir()
		HandleError(err)
		tmpMap[PREV_DIR_ENTRY] = CacheSchema{
			Path:      usrHome,
			Frequency: -1,
		}
		t, err := json.MarshalIndent(tmpMap, "", " ")
		HandleError(err)
		HandleError(os.WriteFile(c.file, t, 077))
	}
}

func (c *Cache) LoadCache() {
	// Load the cache contents to memory
	f, err := os.ReadFile(c.file)
	HandleError(err)
	HandleError(json.Unmarshal(f, &c.contents))
}

func (c *Cache) GetPreviousDir() string {
	return c.contents[PREV_DIR_ENTRY].Path
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

func (c *Cache) GetCacheEntry(entry string) (cacheEntry CacheSchema, ok bool) {
	if cache, found := c.contents[entry]; found {
		return cache, true
	}
	return CacheSchema{}, false
}

func (c *Cache) SetCacheEntry(entry string) {
	// Shoudl check if entry is in cache ,if yes just update
	// else pop from cache and add new entry
	home, _ := os.UserHomeDir()
	if !(entry == home) { //  Do not add home directory to cache
		if cacheEntry, ok := c.GetCacheEntry(filepath.Base(entry)); ok {
			if cacheEntry.Path == entry {
				c.contents[filepath.Base(entry)] = CacheSchema{
					Path:      cacheEntry.Path,
					Frequency: cacheEntry.Frequency + 1,
					LastHit:   time.Now(),
				}
			}
		} else {
			c.contents[filepath.Base(entry)] = CacheSchema{
				Path:      entry,
				Frequency: 0,
				LastHit:   time.Now(),
			}
		}
		if len(c.contents) > (c.maxCap + 1) { // +1 to denote the previous dir store
			c.popCache()
		}
		writeContent, err := json.MarshalIndent(c.contents, "", " ")
		HandleError(err)
		HandleError(os.WriteFile(c.file, writeContent, 077))
	}
}

func (c *Cache) popCache() {
	// Removes according to LFU
	var entryToRemove string
	for i := range c.contents {
		if i == PREV_DIR_ENTRY {
			continue
		}
		if len(entryToRemove) == 0 {
			entryToRemove = i
		} else {
			if c.contents[i].Frequency <= c.contents[entryToRemove].Frequency {
				if c.contents[i].Frequency == c.contents[entryToRemove].Frequency {
					if c.contents[i].LastHit.Before(c.contents[entryToRemove].LastHit) {
						entryToRemove = i
					}
				} else {
					entryToRemove = i
				}
			}
		}
	}
	delete(c.contents, entryToRemove)
}

func (c *Cache) cleanCache() {
	HandleError(os.WriteFile(c.file, []byte{}, 0777))
}
