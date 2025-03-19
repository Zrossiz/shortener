package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zrossiz/shortener/internal/apperrors"
	"github.com/Zrossiz/shortener/internal/dto"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	s Service
}

type Service interface {
	Create(original string) (string, error)
	Get(short string) (string, error)
}

func NewHandler(s Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	var body dto.CreateURL
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	hash, err := h.s.Create(body.OriginalURL)
	if err != nil {
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(rw).Encode(map[string]string{"hash": hash})
	if err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	originalUrl, err := h.s.Get(hash)
	if err != nil {
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, originalUrl, http.StatusMovedPermanently)
}
