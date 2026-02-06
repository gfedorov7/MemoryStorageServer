package cmd

import (
	"MemoryStorageServer/internal/collection"
)

func removeAllExpiredHandler(storage collection.AsyncCollectionInterface) error {
	storage.RemoveAllExpired()
	return nil
}
