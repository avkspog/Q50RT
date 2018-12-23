package main

import (
	"runtime"
	"sync"
	"time"
)

type Cache struct {
	mu          *sync.RWMutex
	items       map[int]Item
	stopCleaner chan struct{}
}

type Item struct {
	Value      interface{}
	Expiration int64
}

const (
	DefaultExpiration      time.Duration = 24 * time.Hour
	DefaultCleanupInterval time.Duration = 5 * time.Second
)

func newCache() *Cache {
	cache := &Cache{
		mu:          &sync.RWMutex{},
		items:       make(map[int]Item),
		stopCleaner: make(chan struct{}),
	}

	go run(cache)
	runtime.SetFinalizer(cache, stop)

	return cache
}

func (c *Cache) Set(key int, object interface{}) {
	c.set(key, object, 0)
}

func (c *Cache) SetExp(key int, object interface{}, expiration time.Duration) {
	c.set(key, object, expiration)
}

func (c *Cache) set(key int, object interface{}, expiration time.Duration) {
	var exp time.Duration

	if expiration == 0 {
		exp = DefaultExpiration
	}

	if expiration > 0 {
		exp = expiration
	}

	c.mu.Lock()
	c.items[key] = Item{
		Value:      object,
		Expiration: time.Now().Add(exp).UnixNano(),
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key int) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	if item.Expired() {
		c.mu.RUnlock()
		return nil, false
	}

	c.mu.RUnlock()
	return item.Value, true
}

func (c *Cache) Delete(key int) {
	c.mu.Lock()
	_, found := c.items[key]
	if found {
		delete(c.items, key)
	}
	c.mu.Unlock()
}

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

func (c *Cache) DeleteExpired() {
	timeNow := time.Now().UnixNano()
	c.mu.Lock()
	for k, item := range c.items {
		if timeNow > item.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

func run(c *Cache) {
	ticker := time.NewTicker(DefaultCleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stopCleaner:
			ticker.Stop()
			return
		}
	}
}

func stop(c *Cache) {
	c.stopCleaner <- struct{}{}
}
