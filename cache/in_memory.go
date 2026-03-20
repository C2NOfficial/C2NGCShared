package cache

import (
	"sync"
)

// Stores the cost of shipment for an order
var ShippingEstimateCache = NewCache[float64] 
// Stores the count for each ip how many times it requested
var IPCache = NewCache[int] 

type entry[T any] struct {
	Value T
}

type Cache[T any] struct {
	entries map[string]entry[T]
	lock    sync.RWMutex //Map is not thread safe
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		entries: make(map[string]entry[T]),
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	e, ok := c.entries[key]
	return e.Value, ok
}

func (c *Cache[T]) Set(key string, value T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries[key] = entry[T]{Value: value}
}
