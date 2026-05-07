package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu           sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
