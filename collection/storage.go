package collection

import (
	"MemoryStorageServer/errors"
	"sync"
)

type AsyncCollection struct {
	collection map[string]MemoryCollection
	mux        sync.RWMutex
}

type AsyncCollectionInterface interface {
	Set(key string, mc MemoryCollection)
	Get(key string) (MemoryCollection, error)
	Remove(key string) bool
}

func NewAsyncCollection() AsyncCollectionInterface {
	return &AsyncCollection{
		collection: make(map[string]MemoryCollection),
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
