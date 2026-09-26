package handlers

import "url-shortener/storage"

type Handler struct {
	store storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{store: store}
}
