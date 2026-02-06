package cmd

import (
	collection2 "MemoryStorageServer/internal/collection"
	"fmt"
	"strconv"
	"time"
)

func setHandler(storage collection2.AsyncCollectionInterface, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("SET command wait 3 arg")
	}

	num, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	memoryCollection, err := collection2.Create(args[1], time.Duration(num)*time.Second, time.Now())
	if err != nil {
		return err
	}

	storage.Set(args[0], memoryCollection)
	return nil
}
