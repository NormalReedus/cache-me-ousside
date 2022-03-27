package cache

import (
	"fmt"
	"log"

	"github.com/NormalReedus/cache-me-ousside/internal/logger"
)

func New(cap uint) *LRUCache {
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
		log.Println(fmt.Errorf("setting the key: %v in the cache should not be done, since that key already exists", key))
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
