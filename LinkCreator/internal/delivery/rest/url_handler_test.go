package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/LinkCreator/creator/internal/domain/dto"
	"github.com/Zrossiz/LinkCreator/creator/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUrlService struct {
	mock.Mock
}

func (m *MockUrlService) Create(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func TestCreateUrl_Success(t *testing.T) {
	mockService := new(MockUrlService)
	handler := NewUrlHandler(mockService)

	testUrl := "https://example.com"
	expectedHash := "abc123"

	mockService.On("Create", testUrl).Return(expectedHash, nil)

	requestBody, _ := json.Marshal(dto.CreateUrl{Url: testUrl})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.CreateUrl(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, response["hash"])

	mockService.AssertExpectations(t)
}

func TestCreateUrl_InvalidMethod(t *testing.T) {
	mockService := new(MockUrlService)
	handler := NewUrlHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/url/", nil)
	rec := httptest.NewRecorder()

	handler.CreateUrl(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid Method")
}

func TestCreateUrl_InvalidPath(t *testing.T) {
	mockService := new(MockUrlService)
	handler := NewUrlHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/invalid/", nil)
	rec := httptest.NewRecorder()

	handler.CreateUrl(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Not Found")
}

func TestCreateUrl_InvalidRequestBody(t *testing.T) {
	mockService := new(MockUrlService)
	handler := NewUrlHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/", bytes.NewBuffer([]byte("invalid json")))
	rec := httptest.NewRecorder()

	handler.CreateUrl(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), apperrors.ErrInvalidRequestBody)
}

func TestCreateUrl_ServiceError(t *testing.T) {
	mockService := new(MockUrlService)
	handler := NewUrlHandler(mockService)

	testUrl := "https://example.com"

	mockService.On("Create", testUrl).Return("", errors.New("internal error"))

	requestBody, _ := json.Marshal(dto.CreateUrl{Url: testUrl})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/url/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.CreateUrl(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), apperrors.ErrInternalServer)

	mockService.AssertExpectations(t)
}
