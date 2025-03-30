package redisdb

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Zrossiz/Redirector/redirector/pkg/apperrors"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisTestSuite struct {
	suite.Suite
	client    *redis.Client
	repo      *RedisRepo
	container testcontainers.Container
}

func (s *RedisTestSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("error starting container: %s", err)
	}

	s.container = container

	ip, err := container.Host(ctx)
	assert.NoError(s.T(), err)

	port, err := container.MappedPort(ctx, "6379")
	assert.NoError(s.T(), err)

	addr := fmt.Sprintf("%s:%s", ip, port.Port())
	s.client, err = Connect(addr, "")
	assert.NoError(s.T(), err)

	s.repo = New(s.client)
}

func (s *RedisTestSuite) TearDownSuite() {
	s.client.Close()
	s.container.Terminate(context.Background())
}

func (s *RedisTestSuite) TestGet_Success() {
	err := s.client.Set("abc123", "https://example.com", 10*time.Minute).Err()
	assert.NoError(s.T(), err)

	url, err := s.repo.Get("abc123")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "https://example.com", url)
}

func (s *RedisTestSuite) TestGet_NotFound() {
	url, err := s.repo.Get("notfound")

	assert.Error(s.T(), err)
	assert.Empty(s.T(), url)
	assert.ErrorIs(s.T(), err, apperrors.ErrNotFound)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
