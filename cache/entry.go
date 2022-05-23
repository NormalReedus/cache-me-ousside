package cache

// newEntry returns a CacheEntry with the given key and data.
func newEntry(key string, data *CacheData) *CacheEntry {
	entry := &CacheEntry{
		key:  key,
		data: data,
	}

	return entry
}

// CacheEntry is used to represent one entry in the LRUCache.
// It is like a node in a linked list.
type CacheEntry struct {
	// key is the name of the entry in the cache.
	// It is usually named after the route that is being cached.
	key string
	// data is an instance of CacheData, which contains both headers and body of an API response.
	data *CacheData
	// next contains a newer CacheEntry in the cache.
	next *CacheEntry
	// prev contains an older CacheEntry in the cache.
	prev *CacheEntry
}

// SetNext will insert a newEntry after the current entry in the linked list.
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

// Key returns the key of the entry.
func (entry CacheEntry) Key() string {
	return entry.key
}

// Data returns the data of the entry encoded in whichever way it was saved in CacheData.
func (entry CacheEntry) Data() *CacheData {
	return entry.data
}

// Prev returns the previous entry in the cache.
func (entry *CacheEntry) Prev() *CacheEntry {
	return entry.prev
}

// Next returns the next entry in the cache.
func (entry *CacheEntry) Next() *CacheEntry {
	return entry.next
}
