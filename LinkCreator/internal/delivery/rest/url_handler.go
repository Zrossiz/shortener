package rest

import (
	"encoding/json"
	"github.com/Zrossiz/LinkCreator/creator/internal/domain"
	"github.com/Zrossiz/LinkCreator/creator/internal/domain/dto"
	"github.com/Zrossiz/LinkCreator/creator/pkg/apperrors"

	"net/http"
)

type UrlHandler struct {
	service domain.UrlService
}

func NewUrlHandler(service domain.UrlService) *UrlHandler {
	return &UrlHandler{service: service}
}

func (h *UrlHandler) CreateUrl(rw http.ResponseWriter, r *http.Request) {
	var body dto.CreateUrl
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	hash, err := h.service.Create(body.Url)
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
