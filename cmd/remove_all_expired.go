package cmd

import (
	"MemoryStorageServer/collection"
)

func removeAllExpiredHandler(storage collection.AsyncCollectionInterface) error {
	storage.RemoveAllExpired()
	return nil
}
