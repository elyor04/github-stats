package cache

import (
	"sync"
	"time"
)

type entry struct {
	data      any
	timestamp time.Time
}

type Cache struct {
	mu  sync.Mutex
	ttl time.Duration
	m   map[string]entry
}

func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, m: map[string]entry{}}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.m[key]
	if !ok || time.Since(e.timestamp) >= c.ttl {
		return nil, false
	}
	return e.data, true
}

func (c *Cache) Set(key string, data any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = entry{data: data, timestamp: time.Now()}
}
