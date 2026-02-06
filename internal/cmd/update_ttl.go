package cmd

import (
	"MemoryStorageServer/internal/collection"
	"fmt"
	"strconv"
	"time"
)

func updateTTLHandler(
	storage collection.AsyncCollectionInterface,
	args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("UPDATE_TTL command wait 2 arg")
	}

	num, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	_, err = storage.UpdateTTL(args[0], time.Duration(num)*time.Second)
	if err != nil {
		return err
	}
	return nil
}
