package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-shortener/dto"
	"url-shortener/internal/storage"
	"url-shortener/models"
)

type Handler struct {
	store   storage.Storage
	baseURL string
}

func NewHandler(store storage.Storage, baseURL string) *Handler {
	return &Handler{store: store, baseURL: baseURL}
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
	newURL, err := requestIsValid(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := h.store.Create(newURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.toShortenResponse(record)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func requestIsValid(w http.ResponseWriter, r *http.Request) (string, error) {
	var req dto.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if !IsValidURL(req.URL) {
		return "", errors.New("invalid url")
	}
	return req.URL, nil
}

func (h *Handler) toShortenResponse(record *models.URLRecord) dto.ShortenResponse {
	return dto.ShortenResponse{
		ID:          record.ID,
		ShortURL:    h.baseURL + "/" + record.ID,
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
		h.StatsHandler(w, r, id)
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

	response := h.toShortenResponse(record)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllLinksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	records, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseRecords := make([]dto.ShortenResponse, 0, len(records))
	for _, record := range records {
		response := h.toShortenResponse(record)
		responseRecords = append(responseRecords, response)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseRecords); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateLink(w http.ResponseWriter, r *http.Request, id string) {
	newURL, err := requestIsValid(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.UpdateById(id, newURL); err != nil {
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

	response := dto.StatsResponse{
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
