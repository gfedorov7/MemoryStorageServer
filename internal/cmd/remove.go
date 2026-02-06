package cmd

import (
	"MemoryStorageServer/internal/collection"
	"fmt"
)

func removeHandler(
	storage collection.AsyncCollectionInterface,
	args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("REMOVE command wait 1 arg")
	}
	isRemoved := storage.Remove(args[0])
	if !isRemoved {
		return fmt.Errorf("cannot use remove for this command")
	}
	return nil
}
