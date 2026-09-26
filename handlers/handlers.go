package handlers

import "url-shortener/storage"

type Handler struct {
	store storage.Storage
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ID          string `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{store: store}
}
