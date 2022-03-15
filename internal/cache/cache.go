package cache

import (
	"container/list"
	"fmt"
)

type LRUCache struct {
	data    map[string][]byte // maps endpoints to JSON data from that endpoint
	queue   list.List         // doubly linked list with strings that represent last recently used endpoints
	maxSize int
}

func New(size int) *LRUCache {
	cache := &LRUCache{
		data:    make(map[string][]byte),
		queue:   *list.New(),
		maxSize: size,
	}

	return cache
}

func (cache *LRUCache) Size() int {
	return len(cache.data)
}

func (cache *LRUCache) Get(key string) ([]byte, bool) { // return a string instead, if that is the format that the router wants to return JSON data
	val, present := cache.data[key]

	// If there was something in cache, move it to MRU
	if present {
		listElement := cache.find(key)

		if listElement == nil {
			fmt.Println(fmt.Errorf("the key %v exists in data map but not in queue", key))

			delete(cache.data, key)  // easiest way to handle this issue is to remove value from cache
			return []byte(""), false // return zero value of map val and as if the key does not exist
		}

		cache.queue.MoveToFront(listElement)
	}

	return val, present
}

func (cache *LRUCache) Set(key string, val *[]byte) {
	cache.data[key] = *val
	cache.queue.PushFront(key)

	// Check if we have exceeded max size and trim LRU if so
	cache.trim()
}

// This is called from cache.Set() after a new key has been added to make sure we don't exceed maxSize
func (cache *LRUCache) trim() {
	// Does not use cache.remove() since remove() uses a key string, whereas trim uses a list element to find target
	if cache.Size() > cache.maxSize {
		lru := cache.queue.Back()

		lruKey := cache.queue.Remove(lru).(string)

		delete(cache.data, lruKey)
	}
}

func (cache *LRUCache) Bust(endpointSelectors []string) {
	// This should be called when any endpoint manipulates the REST API data, i.e. if the request is something that invalidates the cached data
	// this method should be called from the router, and should be passed all the endpoints (or regexes / selectors for endpoints) to remove from the cache with cache.remove()
	// the endpoints (or selectors) are defined in the config that the router loads. This both shows if this specific route is something that should bust the cache and which keys to bust
	// use a selector / regex service (if there are any generics) to convert the 'endpointSelectors' to a slice of strings that are the concrete keys in cache to remove
}

// this is called from cache.Bust() and should just make sure that cache.data and cache.queue stay by removing keys and nodes from both at the same time
func (cache *LRUCache) remove(key string) {
	lruElement := cache.find(key)

	cache.queue.Remove(lruElement)
	delete(cache.data, key)
}

func (cache *LRUCache) find(key string) *list.Element {
	currentElement := cache.queue.Front()

	for {
		if currentElement == nil {
			return nil
		}

		if currentElement.Value == key {
			return currentElement
		}

		currentElement = currentElement.Next()
	}
}
