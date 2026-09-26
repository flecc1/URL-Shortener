package handlers

import "url-shortener/storage"

type Handler struct {
	store storage.MemoryStorage
}

func NewHandler(store storage.MemoryStorage) *Handler {
	return &Handler{store: store}
}
