package cache

import (
	"container/list"
	"encoding/json"
)

type LRUCache struct {
	data    map[string]json.RawMessage // maps endpoints to JSON data from that endpoint
	queue   list.List                  // doubly linked list with strings that represent last recently used endpoints
	maxSize int
}

func New(size int) *LRUCache {
	cache := &LRUCache{
		maxSize: size,
	}

	return cache
}

func (cache *LRUCache) Size() int {
	return len(cache.data)
}

func (cache *LRUCache) Get(endpoint string) json.RawMessage { // return a string instead, if that is the format that the router wants to return JSON data
	// If there is a key for 'endpoint' in cache.data, return the json from the cache and move 'endpoint' to front of queue
	// If there is no key, send a request to endpoint (separate this logic into a service, remember to take method ie GET, POST into account)
	// ... then save the result with cache.set() and return the result

}

func (cache *LRUCache) set(key string, value *json.RawMessage) {
	// This is called from cache.Get()
	// should set key and value in data and then set LRU with cache.queue.PushFront()
	// then call cache.trim() to remove LRU if maxSize has been exceeded
}

func (cache *LRUCache) trim() {
	// This is called from cache.set() after a new key has been added
	// should remove LRU if maxSize has been exceeded by checking cache.Size() - if it is larger than maxSize, remove the LRU from both queue and data, by getting lru with queue.Back() and removing the element with queue.Remove(lru) deleting the cache key with the value returned from queue.Remove(lru)
}

func (cache *LRUCache) Bust(endpointSelectors []string) {
	// This should be called when any endpoint manipulates the REST API data, i.e. if the request is something that invalidates the cached data
	// this method should be called from the router, and should be passed all the endpoints (or regexes / selectors for endpoints) to remove from the cache with cache.remove()
	// the endpoints (or selectors) are defined in the config that the router loads. This both shows if this specific route is something that should bust the cache and which keys to bust
	// use a selector / regex service (if there are any generics) to convert the 'endpointSelectors' to a slice of strings that are the concrete keys in cache to remove
}

func (cache *LRUCache) remove(endpoint string) {
	// this is called from cache.Bust() and should just make sure that cache.data and cache.queue stay by removing keys and nodes from both at the same time
}
