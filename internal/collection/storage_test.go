package collection

import (
	"MemoryStorageServer/internal/errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAsyncCollection_Set_NewKey(t *testing.T) {
	storage := NewAsyncCollection().(*AsyncCollection)

	value := MemoryCollection{
		Value:     []byte("Hello World"),
		CreatedAt: time.Now(),
		TTL:       time.Second * 1,
	}

	storage.Set("key1", value)

	storage.mux.RLock()
	defer storage.mux.RUnlock()

	storedValue, ok := storage.collection["key1"]
	if !ok {
		t.Fatalf("expected key to be present")
	}

	if string(storedValue.Value) != "Hello World" {
		t.Fatalf("unexpected value: %s", storedValue.Value)
	}
}
func TestAsyncCollection_Set_OverwriteValue(t *testing.T) {
	storage := NewAsyncCollection().(*AsyncCollection)

	first := MemoryCollection{
		Value:     []byte("Hello World"),
		CreatedAt: time.Now(),
		TTL:       time.Second * 1,
	}

	second := MemoryCollection{
		Value:     []byte("Hello World2"),
		CreatedAt: time.Now(),
		TTL:       time.Second * 1,
	}

	storage.Set("key", first)
	storage.Set("key", second)

	storage.mux.Lock()
	defer storage.mux.Unlock()

	storedValue, ok := storage.collection["key"]

	if !ok {
		t.Fatalf("expected key to be present")
	}

	if string(storedValue.Value) != "Hello World2" {
		t.Fatalf("value was not overwritten")
	}
}

func TestAsyncCollection_Set_Concurrent(t *testing.T) {
	storage := NewAsyncCollection().(*AsyncCollection)

	const goroutines = 100
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			storage.Set(
				fmt.Sprintf("key-%d", i),
				MemoryCollection{
					Value:     []byte("Hello World"),
					CreatedAt: time.Now(),
					TTL:       time.Second * 1,
				},
			)
		}()
	}

	wg.Wait()

	storage.mux.RLock()
	defer storage.mux.RUnlock()

	if len(storage.collection) != goroutines {
		t.Fatalf("expected %d collections, got %d", goroutines, len(storage.collection))
	}
}

func TestAsyncCollection_Get_ExpiredKey(t *testing.T) {
	storage := NewAsyncCollection().(*AsyncCollection)

	value := MemoryCollection{
		Value:     []byte("Hello World"),
		CreatedAt: time.Now().Add(-1 * time.Hour),
		TTL:       time.Second * 1,
	}

	storage.Set("key", value)

	_, err := storage.Get("key")
	if err == nil {
		t.Fatalf("expected key to be expired")
	}

	if _, ok := err.(errors.ExpiredError); !ok {
		t.Fatalf("error should be of type ExpiredError")
	}
}

func TestAsyncCollection_Get_NotFoundKey(t *testing.T) {
	storage := NewAsyncCollection().(*AsyncCollection)

	_, err := storage.Get("key")

	if err == nil {
		t.Fatalf("expected key should be not found")
	}

	if _, ok := err.(errors.NotFoundError); !ok {
		t.Fatalf("error should be of type NotFoundError")
	}

}
