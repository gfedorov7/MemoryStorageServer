package main

import (
	"MemoryStorageServer/collection"
	"fmt"
	"time"
)

func main() {
	mc, err := collection.Create("test", collection.TypeString, 5)

	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(10 * time.Second)

	asyncCollection := collection.NewAsyncCollection()
	asyncCollection.Set("test", mc)

	asyncCollection.Remove("test")

	test, err := asyncCollection.Get("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(test)
}
