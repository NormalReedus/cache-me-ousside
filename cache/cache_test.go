package cache

import (
	"reflect"
	"testing"
)

var testData CacheData = CacheData{
	Headers: map[string]string{
		"Content-Type": "application/json",
		"X-Test":       "hi-mom",
	},
	Body: []byte(`{"test": "hi mom"}`),
}

func TestRequiredCapacity(t *testing.T) {
	defer func() { recover() }()

	New(0)

	t.Error("Expected cache.New to panic if the capacity is 0")
}

// TODO: split this test into multiple tests to avoid weird interactions (try running a few times to see)
func TestSetGetEntry(t *testing.T) {
	cache := New(2)

	cache.Set("/test1", &testData)

	if cache.Size() != 1 {
		t.Error("Expected cache.Size() to be 1 after adding an entry")
	}

	if cache.Get("/test1") == nil {
		t.Error("Expected cache.Set(\"/test1\") to add an entry to the cache under the key: /test1")
	}

	data := cache.Get("/test1")
	if !reflect.DeepEqual(*data, testData) {
		t.Error("Expected cache.Get to return the exact data that was set with cache.Set")
	}

	cache.Set("/test2", &testData)

	cachedEndpoints := cache.CachedEndpoints()
	if !reflect.DeepEqual(cachedEndpoints, []string{"/test1", "/test2"}) {
		t.Errorf("Expected cache.CachedEndpoints to return the exact keys that were set with cache.Set in the order they were added, it returned: %v", cachedEndpoints)
	}

	// current state - lru = /test1, mru = /test2

	// Should set /test1 to MRU, and implicitly setting /test2 to LRU
	cache.Get("/test1")
	if cache.mru.key != "/test1" {
		t.Error("Expected cache.Get(\"/test1\") to set /test1 to MRU but MRU is: " + cache.mru.key)
	}

	// current state - lru = /test2, mru = /test1

	cache.Set("/test3", &testData)
	if cache.mru.key != "/test3" {
		t.Error("Expected cache.Set(\"/test3\") to set /test3 to MRU but MRU is: " + cache.mru.key)
	}
	// current state - lru = /test2, middle = /test1, mru = /test3

}

func TestEvict(t *testing.T) {
	cache := New(2)

	cache.Set("/test1", &testData)
	cache.Set("/test2", &testData)
	cache.Set("/test3", &testData)

	size := cache.Size()
	if size != 2 {
		t.Errorf("Expected cache size to not go above 2 when capacity is 2, but it is: %d", size)
	}

	if !reflect.DeepEqual(cache.CachedEndpoints(), []string{"/test2", "/test3"}) {
		t.Errorf("Expected cache to evict the oldest entry (/test1), leaving /test2 and /test3 in the cache when going above capacity, but kept: %v", cache.CachedEndpoints())
	}
}

func TestMatchAndBust(t *testing.T) {
	cache := New(5)

	cache.Set("/test1", &testData)
	cache.Set("/test2", &testData)
	cache.Set("/test3", &testData)
	cache.Set("/posts/4", &testData)
	cache.Set("/posts/5", &testData)

	matches := cache.Match([]string{"^/test[1-2]", "^/posts/.+"})
	if reflect.DeepEqual(matches, []string{"/test1", "/test2", "/posts/4", "/posts/5"}) {
		t.Errorf("Expected cache.Match to return 4 matches, but it returned: %v", matches)
	}

	cache.Bust(matches...)
	if cache.CachedEndpoints()[0] != "/test3" {
		t.Errorf("Expected cache.Bust to remove all entries but /test3, but it left: %v", cache.CachedEndpoints())
	}
}
