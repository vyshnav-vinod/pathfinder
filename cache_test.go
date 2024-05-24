// Test for popCache is depended on test for SetCacheEntry!!
// No explicit test is done for GetCacheEntry as it is being
// Checked in two tests and also is a very-less error prone
// function!!

package main

import (
	"os"
	"path/filepath"
	"testing"
)

var c Cache

var (
	toPopEntryLRU = "entry/to/pop/LRU" // will be popped according to LRU
	toPopEntryLFU = "entry/to/pop/LFU" // Will be popped according to LFU

)

func setup() {
	tmpFile, err := os.CreateTemp("", "cache_pf_*.json")
	if err != nil {
		HandleError(err)
	}
	c = Cache{
		file:     tmpFile.Name(),
		maxCap:   10,
		contents: make(map[string]CacheEntry),
	}
}

func teardown() {
	err := os.Remove(c.file)
	if err != nil {
		HandleError(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_CheckCache(t *testing.T) {
	c.CheckCache()
	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		t.Fatalf("CheckCache failed. c.file = %v", c.file)
	}
}

func Test_validateCache(t *testing.T) {
	c.validateCache()
	c.LoadCache()
	want, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Error getting home directory : %v", err)
	}
	got := c.GetPreviousDir()
	if got != want {
		t.Errorf("validateCache() = %s, want = %s", got, want)
	}
}

// TODO: Make this cleaner
func Test_SetCacheEntry(t *testing.T) {
	entries := []pathInfo{
		{
			userInput: "entry",
			restrict:  false,
			path:      "file/path/entry",
		},
		{
			userInput: "path/entry",
			restrict:  true,
			path:      "file/path/entry",
		},
	}
	for _, entry := range entries {
		c.SetCacheEntry(entry)
		_, got := c.GetCacheEntry(entry)
		if !got {
			t.Errorf("SetCacheEntry() failed to add %v to cache", entry)
		}
	}
	test_entries_popCache := []pathInfo{
		{
			userInput: filepath.Base(toPopEntryLRU),
			restrict:  false,
			path:      toPopEntryLRU,
		},
		{
			userInput: filepath.Base(toPopEntryLFU),
			restrict:  false,
			path:      toPopEntryLFU,
		},
	}
	for _, entry := range test_entries_popCache {
		c.SetCacheEntry(entry)
	}
}

// Test for popCache is depended on test for SetCacheEntry
func Test_popCache(t *testing.T) {
	c.LoadCache()
	// fmt.Println(c.contents)
	var popEntries = []pathInfo{
		{
			userInput: filepath.Base(toPopEntryLRU),
			restrict:  false,
			path:      toPopEntryLRU,
		},
		{
			userInput: filepath.Base(toPopEntryLFU),
			restrict:  false,
			path:      toPopEntryLFU,
		},
	}
	for _, e := range popEntries {
		c.popCache()
		if ce, got := c.GetCacheEntry(e); got {
			t.Errorf("popEntry() failed to pop %v ", ce)
		}
	}
}
