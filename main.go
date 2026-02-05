package main

import (
	"MemoryStorageServer/collection"
	"fmt"
	"log"
	"time"
)

func main() {
	col := collection.NewAsyncCollection()
	go col.StartJanitor(time.Duration(5) * time.Second)

	v, err := collection.Create(
		"gleb",
		time.Duration(5)*time.Second,
		time.Now(),
	)
	if err != nil {
		log.Fatal(err)
	}
	col.Set("name", v)

	v, err = col.Get("name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
	b, err := col.UpdateTTL("name", time.Duration(25)*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if b {
		fmt.Println(col.Get("name"))
	}

	time.Sleep(time.Duration(30) * time.Second)

	v, err = col.Get("name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
}
