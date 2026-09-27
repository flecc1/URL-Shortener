package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-shortener/models"
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
type statsResponse struct {
	ID             string `json:"id"`
	CreatedAt      string `json:"created_at"`
	LastAccessedAt string `json:"last_access_at"`
	AccessCount    int    `json:"access_count"`
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
	req := shortenRequest{}

	record, err := h.store.Create(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := toShortenResponse(record)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func requestIsValid(w http.ResponseWriter, r *http.Request) shortenRequest {
	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return shortenRequest{}
	}
	defer r.Body.Close()

	if !IsValidURL(req.URL) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return shortenRequest{}
	}
	return shortenRequest{URL: req.URL}
}

func toShortenResponse(record *models.URLRecord) shortenResponse {
	return shortenResponse{
		ID:          record.ID,
		ShortURL:    "http://localhost:8080/" + record.ID,
		OriginalURL: record.OriginalURL,
		CreatedAt:   record.CreatedAt.Format(time.RFC3339),
	}
}

func (h *Handler) LinkHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/shorten/")
	isStats := strings.HasSuffix(path, "/stats")
	id := strings.TrimSuffix(path, "/stats")

	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if isStats {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.statsHandler(w, r, id)
		return
	}

	switch r.Method {
	case "GET":
		h.getLink(w, r, id)
	case "PUT":
		h.updateLink(w, r, id)
	case "DELETE":
		h.deleteLink(w, r, id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getLink(w http.ResponseWriter, r *http.Request, id string) {
	record, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := toShortenResponse(record)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateLink(w http.ResponseWriter, r *http.Request, id string) {
	req := requestIsValid(w, r)
	if err := h.store.UpdateById(id, req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteLink(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.store.DeleteById(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request, id string) {
	record, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := statsResponse{
		ID:             record.ID,
		AccessCount:    record.AccessedCount,
		CreatedAt:      record.CreatedAt.Format(time.RFC3339),
		LastAccessedAt: record.LastAccessedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	record, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = h.store.IncrementAccess(record.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, record.OriginalURL, http.StatusFound)
}
