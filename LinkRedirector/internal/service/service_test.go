package service

import (
	"errors"
	"testing"

	"github.com/Zrossiz/Redirector/redirector/internal/domain"
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

type mockKafkaProducer struct {
	mock.Mock
}

func (m *mockKafkaProducer) Send(dto domain.UrlKafkaDTO) error {
	args := m.Called(dto)
	return args.Error(0)
}

func TestService_Get_SuccessFromRedis(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	mockKafka := new(mockKafkaProducer)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis, mockKafka)

	mockRedis.On("Get", "abc123").Return("https://example.com", nil)

	dto := domain.GetUrlDTO{OS: "Windows", Short: "abc123", UserIP: "192.168.1.1"}

	url, err := service.Get(dto)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertNotCalled(t, "Get")
	mockKafka.AssertNotCalled(t, "Send")
}

func TestService_Get_SuccessFromPostgres(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	mockKafka := new(mockKafkaProducer)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis, mockKafka)

	mockRedis.On("Get", "abc123").Return("", apperrors.ErrNotFound)
	mockPostgres.On("Get", "abc123").Return("https://example.com", nil)
	mockKafka.On("Send", mock.Anything).Return(nil)

	dto := domain.GetUrlDTO{OS: "Windows", Short: "abc123", UserIP: "192.168.1.1"}

	url, err := service.Get(dto)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
	mockKafka.AssertCalled(t, "Send", mock.Anything)
}

func TestService_Get_ErrorFromRedisAndPostgres(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	mockKafka := new(mockKafkaProducer)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis, mockKafka)

	mockRedis.On("Get", "abc123").Return("", errors.New("redis error"))
	mockPostgres.On("Get", "abc123").Return("", errors.New("postgres error"))

	dto := domain.GetUrlDTO{OS: "Windows", Short: "abc123", UserIP: "192.168.1.1"}

	url, err := service.Get(dto)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Equal(t, "postgres error", err.Error())
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
	mockKafka.AssertNotCalled(t, "Send")
}

func TestService_Get_ErrorFromPostgresOnly(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	mockKafka := new(mockKafkaProducer)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis, mockKafka)

	mockRedis.On("Get", "abc123").Return("", nil)
	mockPostgres.On("Get", "abc123").Return("", errors.New("postgres error"))

	dto := domain.GetUrlDTO{OS: "Windows", Short: "abc123", UserIP: "192.168.1.1"}

	url, err := service.Get(dto)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Equal(t, "postgres error", err.Error())
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
	mockKafka.AssertNotCalled(t, "Send")
}

func TestService_Get_SuccessFromPostgres_WithKafkaError(t *testing.T) {
	mockRedis := new(mockRedisRepo)
	mockPostgres := new(mockPostgresRepo)
	mockKafka := new(mockKafkaProducer)
	logger := zap.NewNop()

	service := NewService(logger, mockPostgres, mockRedis, mockKafka)

	mockRedis.On("Get", "abc123").Return("", apperrors.ErrNotFound)
	mockPostgres.On("Get", "abc123").Return("https://example.com", nil)
	mockKafka.On("Send", mock.Anything).Return(errors.New("kafka error"))

	dto := domain.GetUrlDTO{OS: "Windows", Short: "abc123", UserIP: "192.168.1.1"}

	url, err := service.Get(dto)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	mockRedis.AssertCalled(t, "Get", "abc123")
	mockPostgres.AssertCalled(t, "Get", "abc123")
	mockKafka.AssertCalled(t, "Send", mock.Anything)
}
