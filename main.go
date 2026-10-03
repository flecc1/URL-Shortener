package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("нет .env файла, используются системные переменные")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("не удалось открыть подключение:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("не удалось подключиться к базе:", err)
	}
	fmt.Println("подключение к базе успешно")

	store := storage.NewPostgresStorage(db)
	h := handlers.NewHandler(store, "http://localhost:8080")

	http.HandleFunc("/shorten", h.ShortenHandler)
	http.HandleFunc("/shorten/", h.LinkHandler)
	http.HandleFunc("/", h.RedirectHandler)

	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ошибка запуска сервера:", err)
	}
}
