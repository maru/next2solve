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
	c.mutex.Lock()
	item, ok := c.items[key]
	if !ok {
		c.mutex.Unlock()
		return nil, false
	}

	if item.expiration < time.Now().UnixNano() {
		delete(c.items, key)
		c.mutex.Unlock()
		return nil, false
	}
	c.mutex.Unlock()
	return item.value, true
}

// Set value in item key.
func (c *Cache) Set(key string, value interface{}) {
	expiration := time.Now().Add(c.duration).UnixNano()

	c.mutex.Lock()
	c.items[key] = CacheItem{value, expiration}
	c.mutex.Unlock()
}
