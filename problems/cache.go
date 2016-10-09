// Next problem to solve
// https://github.com/maru/next2solve
//
// Cache for problems
//

package problems

import (
  "errors"
  "time"
  "sync"
)

type Cache struct {
  items map[string]NamespaceCache
  mutex sync.Mutex
}

type NamespaceCache struct {
  items map[string]CacheItem
  Duration  time.Duration
  mutex sync.Mutex
}

type CacheItem struct {
  value interface{}
  expiration int64
}

func NewCache() (*Cache) {
  var cache *Cache
  cache = new(Cache)
  cache.items = make(map[string]NamespaceCache)
  return cache
}

func (c *Cache) CreateNamespace(namespace string) (error) {
  defer c.mutex.Unlock()
  c.mutex.Lock()

  ns, ok := cache.items[namespace]
  if ok {
    return errors.New("Namespace " + namespace + " exists")
  }
  ns.items = make(map[string]CacheItem)
  return nil
}

// Get value from item key, from cache namespace.
// Return value and true if found, otherwise nil and false.
func (c *Cache) Get(namespace string, key string) (interface{}, bool) {
  defer c.mutex.Unlock()
  c.mutex.Lock()

  ns, ok := c.items[namespace]
  if !ok {
    return nil, false
  }

  defer ns.mutex.Unlock()
  ns.mutex.Lock()

  item, ok := ns.items[key]
  if !ok {
    return nil, false
  }
  return item.value, true
}

// Set value in item key, in cache namespace.
// Return error if any.
func (c *Cache) Set(namespace string, key string, value interface{}) (error) {
  defer c.mutex.Unlock()
  c.mutex.Lock()

  ns, ok := c.items[namespace]
  if !ok {
    return errors.New("Namespace " + namespace + " not found.")
  }

  defer ns.mutex.Unlock()
  ns.mutex.Lock()

  expiration := time.Now().Add(ns.Duration).UnixNano()
  ns.items[key] = CacheItem{value, expiration}
  return nil
}
