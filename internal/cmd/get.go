package cmd

import (
	collection2 "MemoryStorageServer/internal/collection"
	"fmt"
)

func getHandler(
	storage collection2.AsyncCollectionInterface,
	args []string) (*collection2.MemoryCollection, error) {

	if len(args) < 1 {
		return nil, fmt.Errorf("GET command wait 1 arg")
	}
	val, err := storage.Get(args[0])
	if err != nil {
		return nil, err
	}
	return &val, nil
}
