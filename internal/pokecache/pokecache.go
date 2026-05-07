package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu           *sync.RWMutex
	interval     time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{cacheEntries: map[string]cacheEntry{}, mu: &sync.RWMutex{}, interval: interval}
	go cache.reapLoop()
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.cacheEntries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	data, ok := cache.cacheEntries[key]
	return data.val, ok
}

func (cache *Cache) reapLoop() {
	ticker := time.NewTicker(cache.interval)
	for tickTimestamp := range ticker.C {
		cache.mu.Lock()
		if len(cache.cacheEntries) >= 1 {
			for key, entry := range cache.cacheEntries {
				elapsed := tickTimestamp.Sub(entry.createdAt)
				if elapsed >= cache.interval {
					delete(cache.cacheEntries, key)
				}
			}
		}
		cache.mu.Unlock()
	}
}
