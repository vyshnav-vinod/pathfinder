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
	toPopEntryLFU = "entry/to/popLFU"  // Will be popped according to LFU

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

func Test_SetCacheEntry(t *testing.T) {
	entry := "random/entry"
	for i := 0; i < 2; i++ {
		c.SetCacheEntry(entry)
		ce, got := c.GetCacheEntry(filepath.Base(entry))
		if !got {
			t.Errorf("SetCacheEntry() failed to add %s to cache", entry)
		}
		if ce.Frequency != i {
			t.Errorf("SetCacheEntry() got frequency = %d, want = %d", ce.Frequency, i)
		}
	}
	c.SetCacheEntry(toPopEntryLRU)
	c.SetCacheEntry(toPopEntryLFU)
}

// Test for popCache is depended on test for SetCacheEntry
func Test_popCache(t *testing.T) {
	c.LoadCache()
	var popEntries = []string{toPopEntryLRU, toPopEntryLFU}
	for _, j := range popEntries {
		c.popCache()
		if entry, ok := c.GetCacheEntry(filepath.Base(j)); ok {
			t.Errorf("popCache() failed to pop %v", entry)
		}
	}
}
