package cmd

import (
	"MemoryStorageServer/collection"
)

func RemoveAllExpiredHandler(storage collection.AsyncCollectionInterface) error {
	storage.RemoveAllExpired()
	return nil
}
