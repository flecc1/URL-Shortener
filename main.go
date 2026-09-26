package main

import (
	"fmt"
	"url-shortener/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	record, err := store.Create("http:/twitch.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("создано: %+v\n", record)
}
