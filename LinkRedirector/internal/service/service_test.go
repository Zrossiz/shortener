package service

import (
	"errors"
	"testing"

	"github.com/Zrossiz/Redirector/redirector/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockRedisRepo struct {
	mock.Mock
}

func (m *mockRedisRepo) Get(hash string) (string, error) {
	args := m.Called(hash)
	return args.String(0), args.Error(1)
}

type mockPostgresRepo struct {
	mock.Mock
}

func (m *mockPostgresRepo) Get(hash string) (string, error) {
	args := m.Called(hash)
	return args.String(0), args.Error(1)
}

func TestService_Get_SuccessFromRedis(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis)

	mockRedis.On("Get", "abc123").Return("https://example.com", nil)

	url, err := service.Get("abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertNotCalled(t, "Get")
}

func TestService_Get_SuccessFromPostgres(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis)

	mockRedis.On("Get", "abc123").Return("", apperrors.ErrNotFound)
	mockPostgres.On("Get", "abc123").Return("https://example.com", nil)

	url, err := service.Get("abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
}

func TestService_Get_ErrorFromRedisAndPostgres(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis)

	mockRedis.On("Get", "abc123").Return("", errors.New("redis error"))
	mockPostgres.On("Get", "abc123").Return("", errors.New("postgres error"))

	url, err := service.Get("abc123")

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Equal(t, "postgres error", err.Error())
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
}

func TestService_Get_ErrorFromPostgresOnly(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis)

	mockRedis.On("Get", "abc123").Return("", nil)
	mockPostgres.On("Get", "abc123").Return("", errors.New("postgres error"))

	url, err := service.Get("abc123")

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Equal(t, "postgres error", err.Error())
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
}
