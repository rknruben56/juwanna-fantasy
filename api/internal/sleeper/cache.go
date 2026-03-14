package sleeper

import (
	"sync"
	"time"
)

// Cache is a thread-safe in-memory cache with TTL support.
type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	data      []byte
	expiresAt time.Time
}

// NewCache creates a new in-memory cache and starts a background cleanup goroutine.
func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]cacheItem),
	}
	go c.cleanup()
	return c
}

// Get retrieves a cached item by key. Returns the data and true if found and not expired.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok || time.Now().After(item.expiresAt) {
		return nil, false
	}
	return item.data, true
}

// Set stores data in the cache with a TTL.
func (c *Cache) Set(key string, data []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		data:      data,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// cleanup periodically removes expired entries.
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
