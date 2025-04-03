package rest

import (
	"strings"

	"github.com/Zrossiz/Redirector/redirector/internal/domain"
	"github.com/Zrossiz/Redirector/redirector/pkg/apperrors"

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

	if r.Method != http.MethodGet {
		http.Error(rw, "Invalid Method", http.StatusMethodNotAllowed)
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	userAgent := r.UserAgent()

	hash := parts[4]

	original, err := h.service.Get(domain.GetUrlDTO{
		OS:     userAgent,
		Short:  hash,
		UserIP: ip,
	})
	if err != nil {
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
	}

	http.Redirect(rw, r, original, http.StatusFound)
}
