package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/Redirector/redirector/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUrlService struct {
	mock.Mock
}

func (m *mockUrlService) Get(hash string) (string, error) {
	args := m.Called(hash)
	return args.String(0), args.Error(1)
}

func TestGetUrl_Success(t *testing.T) {
	mockService := new(mockUrlService)
	handler := NewUrlHandler(mockService)

	mockService.On("Get", "abc123").Return("https://example.com", nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url/abc123", nil)
	w := httptest.NewRecorder()

	handler.GetUrl(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))

	mockService.AssertCalled(t, "Get", "abc123")
}

func TestGetUrl_InvalidURLFormat(t *testing.T) {
	mockService := new(mockUrlService)
	handler := NewUrlHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url", nil)
	w := httptest.NewRecorder()

	handler.GetUrl(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid URL format")

	mockService.AssertNotCalled(t, "Get")
}

func TestGetUrl_InternalServerError(t *testing.T) {
	mockService := new(mockUrlService)
	handler := NewUrlHandler(mockService)

	mockService.On("Get", "abc123").Return("", errors.New(apperrors.ErrInternalServer))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url/abc123", nil)
	w := httptest.NewRecorder()

	handler.GetUrl(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrInternalServer)

	mockService.AssertCalled(t, "Get", "abc123")
}
