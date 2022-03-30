package cache

import (
	"fmt"
	"regexp"

	"github.com/NormalReedus/cache-me-ousside/internal/logger"
	"github.com/NormalReedus/cache-me-ousside/internal/utils"
)

// To use partial memory units, use whole units of lower size instead (e.g. 1.5kb == 1536bytes)
func New(cap uint64) *LRUCache {
	if cap == 0 {
		logger.Panic(fmt.Errorf("cache capacity must be greater than 0"))
	}

	cache := &LRUCache{
		capacity: int(cap),
		entries:  make(map[string]*Entry),
		mru:      nil,
		lru:      nil,
	}

	return cache
}

type LRUCache struct {
	capacity int
	entries  map[string]*Entry
	mru      *Entry
	lru      *Entry
}

func (cache *LRUCache) CachedEndpoints() []string {
	keys := make([]string, 0, len(cache.entries))
	for k := range cache.entries {
		keys = append(keys, k)
	}

	return keys
}

func (cache *LRUCache) Size() int {
	return len(cache.entries)
}

func (cache *LRUCache) Get(key string) *CacheData {
	entry, exists := cache.entries[key]

	if !exists {
		return nil
	}

	// Set as head
	cache.moveToMRU(entry)

	// Unpack data from []byte to CacheData
	data := entry.unmarshalData()

	return &data
}

func (cache *LRUCache) Set(key string, data *CacheData) {
	// You should never set something with a key that already exists
	//... since the cached data should have been returned instead in that case
	if _, exists := cache.entries[key]; exists {
		logger.Warn(fmt.Sprintf("setting the key: %v in the cache should not be done, since that key already exists", key))
		return
	}

	// Ready the data for saving
	entry := newEntry(key, data)

	// If there are no entries, set entry as both head and tail
	if cache.lru == nil && cache.mru == nil {
		cache.entries[key] = cache.setFirst(entry)
	} else {
		// If there are entries, set entry as head
		cache.entries[key] = cache.mru.setNext(entry)
		cache.mru = entry
	}

	// If the cache is full, evict the LRU entry
	if cache.Size() > cache.capacity {
		cache.evict()
	}
}

func (cache *LRUCache) Bust(keys ...string) {
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

func (cache *LRUCache) Match(patterns []string) []string {
	keys := make(utils.Set[string]) // use a set so we don't duplicate keys

	for _, pattern := range patterns {
		patternExp, err := regexp.Compile(pattern)
		if err != nil {
			logger.Error(fmt.Errorf("there was an error finding cache entries with pattern: " + pattern))
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

func (cache *LRUCache) evict() *Entry {
	// Save ref to removed entry
	evicted := cache.lru

	// If there is no lru (cache is empty), don't do anything
	if evicted == nil {
		return nil
	}

	// If there is only one element in the cache
	if cache.Size() == 1 {
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

// Must be used on existing entry to move it to head position
func (cache *LRUCache) moveToMRU(entry *Entry) {
	// If this entry is already head, don't do anything
	if entry == nil || entry == cache.mru {
		return
	}
	if cache.lru == entry {
		cache.lru = entry.next
	}

	cache.mru = entry

	cache.mru.setNext(entry)

}

// If there are no lru or mru, use this to set both to entry
func (cache *LRUCache) setFirst(entry *Entry) *Entry {
	cache.lru = entry
	cache.mru = entry

	entry.prev = nil
	entry.next = nil

	return entry
}
