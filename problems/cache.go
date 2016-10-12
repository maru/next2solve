// Next problem to solve
// https://github.com/maru/next2solve
//
// Cache for problems
//

package problems

import (
	"sync"
	"time"
)

type Cache struct {
	items    map[string]CacheItem
	duration time.Duration
	mutex    sync.Mutex
}

type CacheItem struct {
	value      interface{}
	expiration int64
}

// Return a new cache
func NewCache(duration time.Duration) *Cache {
	var cache *Cache
	cache = new(Cache)
	cache.items = make(map[string]CacheItem)
	cache.duration = time.Duration(duration)
	return cache
}

// Get value from item key.
// Return value and true if found, otherwise nil and false.
func (c *Cache) Get(key string) (interface{}, bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

  if item.expiration < time.Now().UnixNano() {
    delete(c.items, key)
    return nil, false
  }
	return item.value, true
}

// Set value in item key.
// Return error if any.
func (c *Cache) Set(key string, value interface{}) error {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	expiration := time.Now().Add(c.duration).UnixNano()
	c.items[key] = CacheItem{value, expiration}
	return nil
}
