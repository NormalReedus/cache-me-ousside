package main

import (
	"container/list"
	"fmt"
	"log"

	"github.com/fatih/color"
)

type LRUCache struct {
	entries map[string]*list.Element // maps endpoints to JSON data from that endpoint
	queue   *list.List               // doubly linked list with strings that represent last recently used endpoints
	maxSize int
}

func New(size int) *LRUCache {
	cache := &LRUCache{
		entries: make(map[string]*list.Element),
		queue:   list.New(),
		maxSize: size,
	}

	return cache
}

//! DEBUG FJERN IGEN
func (cache *LRUCache) Entries() map[string]*list.Element {
	return cache.entries
}

func (cache *LRUCache) Size() int {
	return len(cache.entries)
}

func (cache *LRUCache) Get(key string) ([]byte, bool) {
	listElement, ok := cache.entries[key]

	// If there is no cached entry, just continue middlewares
	if !ok {
		return []byte(""), false // return zero value of entry val and as if the key does not exist
	}

	// If, for some reason, the cached value is empty
	if listElement == nil || listElement.Value.(cacheElement).Data == nil {
		log.Println(fmt.Errorf("the key %v exists in as an entry but not in queue", key))

		delete(cache.entries, key) // easiest way to handle this issue is to remove value from cache

		return []byte(""), false // return zero value of entry val and as if the key does not exist
	}

	// If there was something in cache, move it to MRU
	cache.queue.MoveToFront(listElement)

	// Return the value in the linked list node (element.Value)
	return *listElement.Value.(cacheElement).Data, ok
}

func (cache *LRUCache) Set(key string, val *[]byte) {
	// Saved data needs both key and data, so we can find map key from the linked list in e.g. evict()
	entry := cacheElement{
		Key:  key,
		Data: val,
	}

	// Save value to both queue and entries
	cache.entries[key] = cache.queue.PushFront(entry)

	// Check if we have exceeded max size and evict LRU if so
	if cache.Size() > cache.maxSize {
		cache.evict()
	}
}

// This is called from cache.Set() after a new key has been added to make sure we don't exceed maxSize
func (cache *LRUCache) evict() {
	// Does not use cache.remove() since remove() uses a key string, whereas evict uses a list element to find target
	lru := cache.queue.Back()

	cacheElm := cache.queue.Remove(lru).(cacheElement)

	delete(cache.entries, cacheElm.Key)

	clr := color.New(color.FgRed, color.Bold)
	log.Printf("%v\n", clr.Sprint("CACHE EVICT: "+cacheElm.Key))
}

func (cache *LRUCache) Bust(endpointSelectors []string) {
	// This should be called when any endpoint manipulates the REST API data, i.e. if the request is something that invalidates the cached data
	// this method should be called from the router, and should be passed all the endpoints (or regexes / selectors for endpoints) to remove from the cache with cache.remove()
	// the endpoints (or selectors) are defined in the config that the router loads. This both shows if this specific route is something that should bust the cache and which keys to bust
	// use a selector / regex service (if there are any generics) to convert the 'endpointSelectors' to a slice of strings that are the concrete keys in cache to remove
}

// this is called from cache.Bust() and should just make sure that cache.data and cache.queue stay synced by removing keys and nodes from both at the same time
func (cache *LRUCache) remove(key string) {
	listElement, ok := cache.entries[key]

	if !ok {
		return
	}

	cache.queue.Remove(listElement)
	delete(cache.entries, key)
}

// content for list.Element.Value
type cacheElement struct {
	Key  string  `json:"key"`  // reference to the key that references this data
	Data *[]byte `json:"data"` // the cached json data
}
