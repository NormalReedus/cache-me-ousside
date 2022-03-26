package cache

func New(cap uint) *LRUCache {
	cache := &LRUCache{
		capacity: int(cap),
		entries:  make(map[string]*Entry),
		head:     nil,
		tail:     nil,
	}

	return cache
}

type LRUCache struct {
	capacity int
	entries  map[string]*Entry
	mru      *Entry
	lru      *Entry
}

func (cache *LRUCache) Size() int {
	return len(cache.entries)
}

func (cache *LRUCache) Get(key string) *Entry {
}

func (cache *LRUCache) Set(key string, entry *Entry) {

	// If there are no entries, set entry as head and tail
	if cache.lru == nil && cache.mru == nil {
		cache.setFirst(entry)
	}

}

func (cache *LRUCache) evict() *Entry {
	if cache.Size() == 0 {
		return nil
	}

	evicted := cache.lru

	if cache.lru == cache.mru {
		// Only one element in the cache

	}

	// if cache.head == nil {
	// 	cache.tail = nil
	// }

	// delete(cache.entries, entry.Key)

	// return entry
}

func (cache *LRUCache) moveToMRU(entry *Entry) *Entry {
	if entry == nil {
		return nil
	}

	// If this entry is head, don't do anything
	if entry == cache.mru {
		return entry
	}

	// If this entry is tail, move to head, but don't link the new tail's prev node
}

// If there are no lru or mru, use this to set both to entry
func (cache *LRUCache) setFirst(entry *Entry) *Entry {
	cache.lru = entry
	cache.mru = entry

	return entry
}

// func (cache *LRUCache) moveToLRU(entry *Entry) {
// 	// if cache.head == nil {
// 	// 	cache.head = entry
// 	// 	cache.tail = entry
// 	// 	return
// 	// }

// 	// if cache.head == cache.tail {
// 	// 	cache.head.Prev = entry
// 	// 	cache.tail.Next = entry
// 	// 	cache.head = entry
// 	// 	cache.tail = entry
// 	// 	return
// 	// }

// 	// entry.Next = cache.head
// 	// cache.head.Prev = entry
// 	// cache.head = entry
// }

type Entry struct {
	key     string
	headers map[string]string
	value   []byte // this is the marshaled json data
	next    *Entry
	prev    *Entry
}

func (entry *Entry) SetNext(newEntry *Entry) *Entry {
	if newEntry == nil {
		return nil
	}

	// If this entry is head (or only entry), set newEntry as head
	if entry.next == nil {
		entry.next = newEntry
		newEntry.prev = entry

		return newEntry
	}

	// If this entry is not head, insert newEntry between this entry and next entry
	nextEntry := entry.next
	entry.next = newEntry
	newEntry.prev = entry
	newEntry.next = nextEntry
	nextEntry.prev = newEntry

	return newEntry
}

func (entry *Entry) SetPrev(newEntry *Entry) *Entry {
	if newEntry == nil {
		return nil
	}

	// If this entry is tail (or only entry), set newEntry as tail
	if entry.prev == nil {
		entry.prev = newEntry
		newEntry.next = entry

		return newEntry
	}

	// If this entry is not tail, insert newEntry between this entry and prev entry
	prevEntry := entry.prev
	entry.prev = newEntry
	newEntry.next = entry
	newEntry.prev = prevEntry
	prevEntry.next = newEntry

	return newEntry
}
