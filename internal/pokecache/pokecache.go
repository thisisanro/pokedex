package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.cache[key]
	if !found {
		return nil, found
	}
	return entry.val, found
}

func (c *Cache) reapLoop() {
	tick := time.NewTicker(c.interval)
	defer tick.Stop()
	for {
		<-tick.C
		c.mu.Lock()
		for key, entry := range c.cache {
			age := time.Since(entry.createdAt)
			if age > c.interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
