package cmd

import (
	"MemoryStorageServer/internal/collection"
	"fmt"
)

func CommandHandler(storage collection.AsyncCollectionInterface, command string,
	args []string) (*collection.MemoryCollection, error) {
	switch command {
	case "GET":
		return getHandler(storage, args)
	case "SET":
		return nil, setHandler(storage, args)
	case "REMOVE":
		return nil, removeHandler(storage, args)
	case "REMOVE_ALL_EXPIRED":
		return nil, removeAllExpiredHandler(storage)
	case "UPDATE_TTL":
		return nil, updateTTLHandler(storage, args)
	default:
		return nil, fmt.Errorf("unknow command")
	}
}
