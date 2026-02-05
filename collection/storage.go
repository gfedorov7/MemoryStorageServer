package collection

import (
	"MemoryStorageServer/errors"
	"sync"
	"time"
)

type AsyncCollection struct {
	collection map[string]MemoryCollection
	mux        sync.RWMutex
	stop       chan struct{}
	Clock      Clock
}

type AsyncCollectionInterface interface {
	Set(key string, mc MemoryCollection)
	Get(key string) (MemoryCollection, error)
	Remove(key string) bool
	RemoveAllExpired()
	UpdateTTL(key string, ttl time.Duration) (bool, error)
	StartJanitor(interval time.Duration)
	StopJanitor()
}

func NewAsyncCollection() AsyncCollectionInterface {
	return &AsyncCollection{
		collection: make(map[string]MemoryCollection),
		stop:       make(chan struct{}),
		Clock:      RealClock{},
	}
}

func (c *AsyncCollection) Set(key string, value MemoryCollection) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.collection[key] = value
}

func (c *AsyncCollection) Get(key string) (MemoryCollection, error) {
	c.mux.RLock()
	v, ok := c.collection[key]
	c.mux.RUnlock()

	if !ok {
		return MemoryCollection{}, errors.NotFoundError{Arg: key}
	} else if v.IsExpired(c.Clock.Now()) {
		c.mux.Lock()
		delete(c.collection, key)
		c.mux.Unlock()
		return MemoryCollection{}, errors.ExpiredError{Arg: key}
	}

	copyValue := make([]byte, len(v.Value))
	copy(copyValue, v.Value)
	v.Value = copyValue
	return v, nil
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

func (c *AsyncCollection) RemoveAllExpired() {
	var deleteKeys []string
	c.mux.RLock()
	for key, value := range c.collection {
		if value.IsExpired(c.Clock.Now()) {
			deleteKeys = append(deleteKeys, key)
		}
	}
	c.mux.RUnlock()

	c.mux.Lock()
	for _, key := range deleteKeys {
		delete(c.collection, key)
	}
	c.mux.Unlock()
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
	value.CreatedAt = c.Clock.Now()
	c.collection[key] = value
	return true, nil
}

func (c *AsyncCollection) StartJanitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.RemoveAllExpired()
			case <-c.stop:
				return
			}
		}
	}()
}

func (c *AsyncCollection) StopJanitor() {
	close(c.stop)
}
