package main

import (
	"fmt"
	"net/http"
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	h := handlers.NewHandler(store, "http://localhost:8080")

	http.HandleFunc("/shorten", h.ShortenHandler)
	http.HandleFunc("/shorten/", h.LinkHandler)
	http.HandleFunc("/", h.RedirectHandler)
	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ошибка запуска сервера:", err)
	}
}
