package cache

import (
	"log"
	"sync"
	"time"
)

// Cache represents a simple in-memory cache.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	value      []byte
	expiration time.Time
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]cacheEntry),
	}
}

// Set adds a value to the cache with a specified duration.
func (c *Cache) Set(key string, value []byte, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(duration),
	}
	log.Printf("Set cache for key: %s", key)
}

// Get retrieves a value from the cache. It returns nil if the key does not exist or has expired.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.entries[key]
	if !found || time.Now().After(entry.expiration) {
		return nil, false
	}
	return entry.value, true
}

// ClearCashe removes all entries from the cache.
func (c *Cache) ClearCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cacheEntry)
	log.Println("Cleared all cache")
}
