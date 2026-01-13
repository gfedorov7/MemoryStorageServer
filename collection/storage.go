package collection

import (
	"MemoryStorageServer/errors"
	"sync"
	"time"
)

type AsyncCollection struct {
	collection map[string]MemoryCollection
	mux        sync.RWMutex
}

type AsyncCollectionInterface interface {
	Set(key string, mc MemoryCollection)
	Get(key string) (MemoryCollection, error)
	Remove(key string) bool
	RemoveAllExpired()
	UpdateTTL(key string, ttl time.Duration) (bool, error)
	StartJanitor(interval time.Duration)
}

func NewAsyncCollection() AsyncCollectionInterface {
	return &AsyncCollection{
		collection: make(map[string]MemoryCollection),
	}
}

func (c *AsyncCollection) UpdateTTL(key string, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		return false, errors.TTLError{}
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	value, ok := c.collection[key]
	if !ok {
		return false, errors.NotFoundError{}
	}

	value.TTL = ttl
	value.CreatedAt = getCurrentTime()
	c.collection[key] = value
	return true, nil
}

func (c *AsyncCollection) RemoveAllExpired() {
	c.mux.Lock()
	defer c.mux.Unlock()
	for key, value := range c.collection {
		if value.IsExpired() {
			delete(c.collection, key)
		}
	}
}

func (c *AsyncCollection) Set(key string, value MemoryCollection) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.collection[key] = value
}

func (c *AsyncCollection) Get(key string) (MemoryCollection, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	v, ok := c.collection[key]

	if ok && v.IsExpired() {
		c.mux.RUnlock()
		c.mux.Lock()
		delete(c.collection, key)
		c.mux.Unlock()
		return MemoryCollection{}, errors.ExpiredError{Arg: key}
	} else if ok {
		return v, nil
	} else {
		return MemoryCollection{}, errors.NotFoundError{Arg: key}
	}
}

func (c *AsyncCollection) Remove(key string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.collection[key]; ok {
		delete(c.collection, key)
		return true
	}
	return false
}

func (c *AsyncCollection) StartJanitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.RemoveAllExpired()
		}
	}()
}
