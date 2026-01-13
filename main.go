package main

import (
	"MemoryStorageServer/collection"
	"fmt"
	"time"
)

func main() {
	col := collection.NewAsyncCollection()

	mc, _ := collection.Create("test", collection.TypeString, time.Duration(10)*time.Second)
	col.Set("test", mc)

	v, _ := col.Get("test")
	fmt.Println(v)

	time.Sleep(15 * time.Second)
	col.StartJanitor(time.Duration(10) * time.Second)
	time.Sleep(15 * time.Second)

	_, err := col.UpdateTTL("test", time.Duration(200)*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	v, _ = col.Get("test")
	fmt.Println(v)
}
