package cache

import (
	"encoding/json"
	"fmt"

	"github.com/NormalReedus/cache-me-ousside/internal/logger"
)

func newEntry(key string, data *CacheData) *CacheEntry {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error(fmt.Errorf("there was an error creating the entry: %s", key))
		return nil
	}

	entry := &CacheEntry{
		key:  key,
		data: jsonData,
	}

	return entry
}

type CacheEntry struct {
	key  string
	data []byte // marshaled json of CacheData
	next *CacheEntry
	prev *CacheEntry
}

func (entry *CacheEntry) SetNext(newEntry *CacheEntry) *CacheEntry {
	if newEntry == nil {
		return nil
	}

	// If this entry is not head, insert newEntry between this entry and next entry
	nextEntry := entry.next
	entry.next = newEntry
	newEntry.prev = entry
	// if entry is head, nextEntry is nil and will correctly set the newEntry's next as nil, because newEntry should become head
	newEntry.next = nextEntry
	// If nextEntry is not nil (entry is not head), set the prev of nextEntry to newEntry
	if nextEntry != nil {
		nextEntry.prev = newEntry
	}

	return newEntry
}

func (entry *CacheEntry) UnmarshalData() CacheData {
	return NewCacheDataFromJSON(entry.data)
}

func (entry CacheEntry) Key() string {
	return entry.key
}

func (entry CacheEntry) Data() []byte {
	return entry.data
}

func (entry *CacheEntry) Prev() *CacheEntry {
	return entry.prev
}

func (entry *CacheEntry) Next() *CacheEntry {
	return entry.next
}
