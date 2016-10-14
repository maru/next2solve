// Next problem to solve
// https://github.com/maru/next2solve
//
// Tests for cache.go functionality
//
package problems

import (
	"testing"
	"time"
)

//
const (
	cacheKey   = "key"
	cacheValue = "value"
)

// Test create a new cache.
func TestNewCache(t *testing.T) {
	cache := NewCache(time.Minute)
	if _, ok := cache.Get(cacheKey); ok {
		t.Fatalf("Expected empty cache")
	}
}

// Test create a new cache.
func TestSetGetCache(t *testing.T) {
	duration := time.Second
	cache := NewCache(duration)
	cache.Set(cacheKey, cacheValue)
	time.Sleep(duration / 2)
	obj, ok := cache.Get(cacheKey)
	if !ok {
		t.Fatalf("Expected object found")
	}
	if obj.(string) != cacheValue {
		t.Fatalf("Expected value %s", cacheValue)
	}
	time.Sleep(duration / 2)
	obj, ok = cache.Get(cacheKey)
	if ok {
		t.Fatalf("Expected object not found")
	}
	if obj != nil {
		t.Fatalf("Expected nil object")
	}
}
