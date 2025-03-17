package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Create(rw http.ResponseWriter, r *http.Request)
	Get(rw http.ResponseWriter, r *http.Request)
}

func NewRouter(h Handler) http.Handler {
	r := chi.NewRouter()

	r.Post("/api/v1/url", h.Create)
	r.Get("/api/v1/url/{hash}", h.Get)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
