package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockPostgresRepo struct {
	mock.Mock
}

func (m *MockPostgresRepo) Create(original, hash string) error {
	args := m.Called(original, hash)
	return args.Error(0)
}

type MockRedisRepo struct {
	mock.Mock
}

func (m *MockRedisRepo) Create(hash, original string) error {
	args := m.Called(hash, original)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	log := zap.NewNop()
	postgresMock := new(MockPostgresRepo)
	redisMock := new(MockRedisRepo)
	svc := NewService(log, postgresMock, redisMock)

	originalURL := "https://example.com"
	expectedHash := svc.generateHash(originalURL)

	postgresMock.On("Create", originalURL, expectedHash).Return(nil)
	redisMock.On("Create", expectedHash, originalURL).Return(nil)

	hash, err := svc.Create(originalURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)

	postgresMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestCreate_PostgresError(t *testing.T) {
	log := zap.NewNop()
	postgresMock := new(MockPostgresRepo)
	redisMock := new(MockRedisRepo)
	svc := NewService(log, postgresMock, redisMock)

	originalURL := "https://example.com"
	expectedHash := svc.generateHash(originalURL)

	postgresMock.On("Create", originalURL, expectedHash).Return(errors.New("db error"))

	hash, err := svc.Create(originalURL)

	assert.Error(t, err)
	assert.Empty(t, hash)

	postgresMock.AssertExpectations(t)
	redisMock.AssertNotCalled(t, "Create")
}

func TestCreate_RedisError(t *testing.T) {
	log := zap.NewNop()
	postgresMock := new(MockPostgresRepo)
	redisMock := new(MockRedisRepo)
	svc := NewService(log, postgresMock, redisMock)

	originalURL := "https://example.com"
	expectedHash := svc.generateHash(originalURL)

	postgresMock.On("Create", originalURL, expectedHash).Return(nil)
	redisMock.On("Create", expectedHash, originalURL).Return(errors.New("cache error"))

	hash, err := svc.Create(originalURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)

	postgresMock.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestGenerateHash(t *testing.T) {
	log := zap.NewNop()
	postgresMock := new(MockPostgresRepo)
	redisMock := new(MockRedisRepo)
	svc := NewService(log, postgresMock, redisMock)

	url1 := "https://example.com"
	url2 := "https://another.com"

	hash1 := svc.generateHash(url1)
	hash2 := svc.generateHash(url2)

	assert.Len(t, hash1, 7)
	assert.Len(t, hash2, 7)
	assert.NotEqual(t, hash1, hash2)

	assert.Equal(t, hash1, svc.generateHash(url1))
	assert.Equal(t, hash2, svc.generateHash(url2))
}
