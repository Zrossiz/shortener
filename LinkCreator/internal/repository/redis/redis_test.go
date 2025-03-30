package redisdb

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	redisAddr string
}

func (s *RedisTestSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(s.T(), err)

	s.container = container

	ip, err := container.Host(ctx)
	assert.NoError(s.T(), err)

	port, err := container.MappedPort(ctx, "6379")
	assert.NoError(s.T(), err)

	s.redisAddr = fmt.Sprintf("%s:%s", ip, port.Port())

	client, err := Connect(s.redisAddr, "")
	assert.NoError(s.T(), err)

	s.client = client
	s.repo = New(client)
}

func (s *RedisTestSuite) TearDownSuite() {
	s.client.Close()
	s.container.Terminate(context.Background())
}

func (s *RedisTestSuite) TestCreate_Success() {
	hash := "abc123"
	original := "https://example.com"
	err := s.repo.Create(hash, original)
	assert.NoError(s.T(), err)

	val, err := s.client.Get(hash).Result()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), original, val)
}

func (s *RedisTestSuite) TestCreate_Error() {
	hash := ""
	original := ""

	err := s.repo.Create(hash, original)
	assert.Error(s.T(), err)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
