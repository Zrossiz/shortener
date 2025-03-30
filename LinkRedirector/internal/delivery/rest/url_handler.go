package rest

import (
	"encoding/json"
	"github.com/Zrossiz/LinkRedirector/redirector/internal/domain"
	"github.com/Zrossiz/LinkRedirector/redirector/pkg/apperrors"
	"strings"

	"net/http"
)

type UrlHandler struct {
	service domain.UrlService
}

func NewUrlHandler(service domain.UrlService) *UrlHandler {
	return &UrlHandler{service: service}
}

func (h *UrlHandler) GetUrl(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) < 5 {
		http.Error(rw, "Invalid URL format", http.StatusBadRequest)
		return
	}

	hash := parts[4]

	original, err := h.service.Get(hash)
	if err != nil {

	}
}
