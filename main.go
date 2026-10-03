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
	"url-shortener/internal/midleware"
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

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.ShortenHandler)
	mux.HandleFunc("GET /shorten/all", h.GetAllHandler)
	mux.HandleFunc("GET /shorten/{id}", h.GetLink)
	mux.HandleFunc("PUT /shorten/{id}", h.UpdateLink)
	mux.HandleFunc("DELETE /shorten/{id}", h.DeleteLink)
	mux.HandleFunc("GET /shorten/{id}/stats", h.StatsHandler)
	mux.HandleFunc("GET /{id}", h.RedirectHandler)

	fmt.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", midleware.LoggingMiddleware(mux)); err != nil {
		fmt.Println("ошибка запуска сервера:", err)
	}
}
