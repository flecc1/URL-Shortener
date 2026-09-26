package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"url-shortener/storage"
)

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

func IsValidURL(URL string) bool {
	parsed, err := url.ParseRequestURI(URL)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	if parsed.Host == "" {
		return false
	}
	return true
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !IsValidURL(req.URL) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := h.store.Create(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := shortenResponse{
		ID:          record.ID,
		ShortURL:    "http://localhost:8080/" + record.ID,
		OriginalURL: req.URL,
		CreatedAt:   record.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
