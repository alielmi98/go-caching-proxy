package cache

import (
	"log"
	"sync"
	"time"
)

// Cache represents a simple in-memory cache with size limits and cleanup.
type Cache struct {
	mu            sync.RWMutex
	entries       map[string]cacheEntry
	maxEntries    int
	cleanupTicker *time.Ticker
	stopCleanup   chan bool
}

type cacheEntry struct {
	value        []byte
	expiration   time.Time
	lastAccessed time.Time
	size         int
}

//create a new Cache instance with default values.
func NewCache() *Cache {
	return NewCacheWithConfig(1000, 5*time.Minute) // Default: 1000 entries, 5min cleanup
}

// create new Cache instance with configurable limits.
func NewCacheWithConfig(maxEntries int, cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		entries:     make(map[string]cacheEntry),
		maxEntries:  maxEntries,
		stopCleanup: make(chan bool),
	}

	//  periodic cleanup
	if cleanupInterval > 0 {
		cache.cleanupTicker = time.NewTicker(cleanupInterval)
		go cache.periodicCleanup()
	}

	return cache
}

// Set adds a value to the cache with a specified duration.
func (c *Cache) Set(key string, value []byte, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict entries due to size limit
	if len(c.entries) >= c.maxEntries && len(c.entries) > 0 {
		c.evictOldest()
	}

	c.entries[key] = cacheEntry{
		value:        value,
		expiration:   time.Now().Add(duration),
		lastAccessed: time.Now(),
		size:         len(value),
	}
	log.Printf("Set cache for key: %s (size: %d bytes)", key, len(value))
}

// Get retrieves a value from the cache. It returns nil if the key does not exist or has expired.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, found := c.entries[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	// Check if entry has expired
	if time.Now().After(entry.expiration) {
		// Remove any expired entry
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		log.Printf("Removed expired cache entry: %s", key)
		return nil, false
	}

	// Update last accessed time
	c.mu.Lock()
	entry.lastAccessed = time.Now()
	c.entries[key] = entry
	c.mu.Unlock()

	return entry.value, true
}

// ClearCache removes all entries from the cache.
func (c *Cache) ClearCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cacheEntry)
	log.Println("Cleared all cache")
}

// remove expired entries
func (c *Cache) periodicCleanup() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.cleanupExpired()
		case <-c.stopCleanup:
			c.cleanupTicker.Stop()
			return
		}
	}
}

//remove all expired entries from the cache
func (c *Cache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	removedCount := 0

	for key, entry := range c.entries {
		if now.After(entry.expiration) {
			delete(c.entries, key)
			removedCount++
		}
	}

	if removedCount > 0 {
		log.Printf("Cleanup removed %d expired cache entries", removedCount)
	}
}

// remove the oldest accessed entry when cache is full
func (c *Cache) evictOldest() {
	if len(c.entries) == 0 {
		return
	}

	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range c.entries {
		if first {
			oldestKey = key
			oldestTime = entry.lastAccessed
			first = false
			continue
		}
		if entry.lastAccessed.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccessed
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
		log.Printf("Evicted oldest cache entry: %s (last accessed: %v)", oldestKey, oldestTime)
	}
}

//return cach statistics
func (c *Cache) GetStats() (currentSize, maxSize int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries), c.maxEntries
}

// stop the periodic cleanup goroutine
func (c *Cache) StopCleanup() {
	if c.stopCleanup != nil {
		c.stopCleanup <- true
	}
}