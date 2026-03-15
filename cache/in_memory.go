package cache

import (
	"sync"
	"time"
)

var ShippingEstimateCache = NewCache[float64](30 * time.Minute)
var IPCache = NewCache[int](30 * time.Minute)

type entry[T any] struct {
	Value     T
	ExpiresAt time.Time
}

type Cache[T any] struct {
	entries map[string]entry[T]
	ttl     time.Duration
	lock    sync.RWMutex //map is not thread safe in go.
}

func NewCache[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		entries: make(map[string]entry[T]),
		ttl:     ttl,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().UTC().After(e.ExpiresAt) {
		var zero T
		return zero, false
	}
	return e.Value, true
}

func (c *Cache[T]) Set(key string, value T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries[key] = entry[T]{
		Value:     value,
		ExpiresAt: time.Now().UTC().Add(c.ttl),
	}
}
