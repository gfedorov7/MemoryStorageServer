package cmd

import (
	"MemoryStorageServer/internal/collection"
	"fmt"
	"strconv"
	"time"
)

func setHandler(storage collection.AsyncCollectionInterface, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("SET command wait 3 arg")
	}

	num, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	memoryCollection, err := collection.Create(args[1], time.Duration(num)*time.Second, collection.RealClock{}.Now())
	if err != nil {
		return err
	}

	storage.Set(args[0], memoryCollection)
	return nil
}
