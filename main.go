package main

import (
	"MemoryStorageServer/collection"
	"fmt"
)

func main() {
	col := collection.NewAsyncCollection()

	mc, _ := collection.Create("test", collection.TypeString, 10)
	col.Set("test", mc)

	v, _ := col.Get("test")
	fmt.Println(v)

	_, err := col.UpdateTTL("test", 200)
	if err != nil {
		fmt.Println(err)
	}

	v, _ = col.Get("test")
	fmt.Println(v)
}
