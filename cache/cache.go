// Package cache provides a simple LRU cache featuring coupled linked list
// and map data structures to allow for easy lookups and ordered entries.
package cache

import (
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/magnus-bb/cache-me-ousside/internal/logger"
)

// New returns an LRUCache with the given capacity and optionally a unit to use memory-based cache limit.
// To use partial memory units, use whole units of lower size instead (e.g. 1.5kb == 1536b).
// TODO: handle the passed cap unit
func New(capacity uint64, capacityUnit string) (*LRUCache, error) {
	if capacity == 0 {
		return nil, errors.New("cache capacity must be greater than 0")
	}

	cache := &LRUCache{
		capacity: int(capacity),
		entries:  make(map[string]*CacheEntry),
		mru:      nil,
		lru:      nil,
	}

	return cache, nil
}

// LRUCache represents all entries in the cache, it's capacity limit, and the first and last entries.
type LRUCache struct {
	mutex    sync.RWMutex
	capacity int
	entries  map[string]*CacheEntry
	mru      *CacheEntry
	lru      *CacheEntry
}

// CachedKeys returns a slice of the keys of all cached entries.
// NOTE: Does not always return keys in the order they were added.
func (cache *LRUCache) CachedKeys() []string {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	keys := make([]string, 0, len(cache.entries))
	for k := range cache.entries {
		keys = append(keys, k)
	}

	return keys
}

// Size returns the number of entries currently saved in the cache.
func (cache *LRUCache) Size() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	return len(cache.entries)
}

// Get returns the CacheData of the entry saved under the given key.
func (cache *LRUCache) Get(key string) *CacheData {
	// Write lock is used since Get will also rearrange the order of entries
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, exists := cache.entries[key]

	if !exists {
		return nil
	}

	// Set fetched entry as head
	cache.moveToMRU(entry)

	return entry.Data()
}

// Set saves an entry with the given CacheData under the given key in the cache.
func (cache *LRUCache) Set(key string, data *CacheData) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// You should never set something with a key that already exists
	//... since the cached data should have been returned instead in that case
	if _, exists := cache.entries[key]; exists {
		logger.Warn(fmt.Sprintf("the key: %q already exists in the cache and has been ignored", key))
		return
	}

	// Ready the data for saving
	entry := newEntry(key, data)

	// If there are no entries, set entry as both head and tail
	if cache.lru == nil && cache.mru == nil {
		cache.entries[key] = cache.setFirst(entry)
	} else {
		// If there are entries, set entry as head
		cache.entries[key] = cache.mru.SetNext(entry)
		cache.mru = entry
	}

	// If the cache is full, evict the LRU entry
	if len(cache.entries) > cache.capacity { // we don't use Size, since that has its own lock
		cache.evictLRU()
	}
}

// Bust will remove all entries saved under the given keys from the cache.
func (cache *LRUCache) Bust(keys ...string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for _, entryKey := range keys {
		entry, exists := cache.entries[entryKey]

		// It's completely fine if this thing doesn't exist
		if !exists {
			continue
		}

		delete(cache.entries, entryKey)

		if entry == cache.lru {
			cache.lru = entry.next
			if cache.lru != nil {
				cache.lru.prev = nil
			}
		}

		if entry == cache.mru {
			cache.mru = entry.prev
			if cache.mru != nil {
				cache.mru.next = nil
			}
		}

		// It might look weird to not stop execution here if entry were both mru and lru, but if both of those
		// were true, both of these conditions would be nil and thus skipped
		if entry.prev != nil {
			entry.prev.next = entry.next
		}

		if entry.next != nil {
			entry.next.prev = entry.prev
		}

		logger.CacheBust(entryKey)
	}
}

// Match returns a slice of keys of the entries in the cache that match the given patterns.
// The patterns are hydrated with URL parameters from paramMap before being compiled as regex.
// If an empty slice of patterns is passed, all keys are returned (matching everything).
func (cache *LRUCache) Match(patterns []string, paramMap map[string]string) []string {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	keys := make(Set[string]) // use a set so we don't duplicate keys

	if len(patterns) == 0 {
		// If empty slice of patterns is passed, return all keys (match all)
		patterns = []string{"."}

	} else {
		// If there are any route params (/:id for example), insert the actual values into the pattern before compiling regex
		patterns = HydrateParams(paramMap, patterns)
	}

	for _, pattern := range patterns {
		patternExp, err := regexp.Compile(pattern)
		if err != nil {
			logger.Error(fmt.Errorf("there was an error finding cache entries with RegExp pattern: %q", pattern))
			continue
		}

		for key := range cache.entries {
			if patternExp.MatchString(key) {
				keys.Add(key)
			}
		}
	}

	return keys.Elements()
}

// evictLRU removes the least recently used entry from the cache to make room for new entries.
func (cache *LRUCache) evictLRU() *CacheEntry {
	// No need to lock mutex here, this is not an atomic operation
	// which means that the calling operations (Set) will lock the mutex

	// Save ref to removed entry
	evicted := cache.lru

	// If there is no lru (cache is empty), don't do anything
	if evicted == nil {
		return nil
	}

	// If there is only one element in the cache
	if len(cache.entries) == 1 { // Don't use Size since calling operation will handle mutex lock
		cache.lru = nil
		cache.mru = nil

		return nil

	} else {
		// If there is more than one element in the cache
		// Point tail to second-to-last entry
		cache.lru = evicted.next

		// Dereference the removed entry by pointing the prev entry of second-to-last entry to nil
		evicted.next.prev = nil
	}

	// Remove entry from map
	delete(cache.entries, evicted.key)

	logger.CacheEvict(evicted.key)

	return evicted
}

// moveToMRU moves the given entry to the most recently used position in the cache.
// NOTE: Must be used on existing entry, cannot be used to add new entries.
func (cache *LRUCache) moveToMRU(entry *CacheEntry) {
	// No need to lock mutex here, this is not an atomic operation
	// which means that the calling operations (Get) will lock the mutex

	// If this entry is already head (or only entry), don't do anything
	if entry == nil || entry == cache.mru {
		return
	}

	if cache.lru == entry {
		cache.lru = entry.next
		cache.lru.prev = nil
	}

	cache.mru.SetNext(entry)

	cache.mru = entry
}

// setFirst sets the given entry as the first entry in the cache.
// This is used because it takes som special setup to add the first node to a linked list.
func (cache *LRUCache) setFirst(entry *CacheEntry) *CacheEntry {
	// No need to lock mutex here, this is not an atomic operation
	// which means that the calling operations (Set) will lock the mutex

	cache.lru = entry
	cache.mru = entry

	entry.prev = nil
	entry.next = nil

	return entry
}

func (cache *LRUCache) Entries() map[string]*CacheEntry {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	return cache.entries
}

func (cache *LRUCache) MRU() *CacheEntry {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	return cache.mru
}

func (cache *LRUCache) LRU() *CacheEntry {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	return cache.lru
}
