package handlers

import (
	"net/http"

	"zametki/internal/storage"
)

type Handler struct {
	db *storage.Postgres
}

func New(db *storage.Postgres) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.health)

	// и со слешем, и без — чтобы /users/1 работало
	mux.HandleFunc("/users", h.users)
	mux.HandleFunc("/users/", h.users)

	return mux
}
