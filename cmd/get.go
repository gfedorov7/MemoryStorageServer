package cmd

import (
	"MemoryStorageServer/collection"
	"fmt"
)

func getHandler(
	storage collection.AsyncCollectionInterface,
	args []string) (*collection.MemoryCollection, error) {

	if len(args) < 1 {
		return nil, fmt.Errorf("GET command wait 1 arg")
	}
	val, err := storage.Get(args[0])
	if err != nil {
		return nil, err
	}
	return &val, nil
}
