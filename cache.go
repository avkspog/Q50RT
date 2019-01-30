package main

import (
	"runtime"
	"sync"
	"time"
)

type Cache struct {
	mu          *sync.RWMutex
	Items       map[int]Item
	stopCleaner chan struct{}
}

type Item struct {
	Value      interface{}
	Expiration int64
}

const (
	DefaultExpiration      = 24 * time.Hour
	DefaultCleanupInterval = 5 * time.Second
)

func NewCache() *Cache {
	cache := &Cache{
		mu:          &sync.RWMutex{},
		Items:       make(map[int]Item),
		stopCleaner: make(chan struct{}),
	}

	go runCleaner(cache)
	runtime.SetFinalizer(cache, stopCleaner)

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
	c.Items[key] = Item{
		Value:      object,
		Expiration: time.Now().Add(exp).UnixNano(),
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key int) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.Items[key]
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
	_, found := c.Items[key]
	if found {
		delete(c.Items, key)
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
	for k, item := range c.Items {
		if timeNow > item.Expiration {
			delete(c.Items, k)
		}
	}
	c.mu.Unlock()
}

func runCleaner(c *Cache) {
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

func stopCleaner(c *Cache) {
	c.stopCleaner <- struct{}{}
}
