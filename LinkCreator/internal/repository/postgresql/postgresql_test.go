package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestSuite struct {
	suite.Suite
	db        *sql.DB
	repo      *PosgresRepo
	dbURI     string
	container testcontainers.Container
}

func (s *PostgresTestSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("start container error: %s", err)
	}

	s.container = container

	ip, err := container.Host(ctx)
	assert.NoError(s.T(), err)

	port, err := container.MappedPort(ctx, "5432")
	assert.NoError(s.T(), err)

	s.dbURI = fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", ip, port.Port())

	s.db, err = Connect(s.dbURI)
	assert.NoError(s.T(), err)

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		original TEXT NOT NULL UNIQUE,
		short VARCHAR(7) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`)
	assert.NoError(s.T(), err)

	s.repo = NewPostgresRepo(s.db)
}

func (s *PostgresTestSuite) TearDownSuite() {
	s.db.Close()
	s.container.Terminate(context.Background())
}

func (s *PostgresTestSuite) TestCreate_Success() {
	err := s.repo.Create("https://example1.com", "abc1234")
	assert.NoError(s.T(), err)
}

func (s *PostgresTestSuite) TestCreate_Duplicate() {
	url := "https://example.com"
	hash := "abc123"

	err := s.repo.Create(url, hash)
	assert.NoError(s.T(), err)

	err = s.repo.Create(url, hash)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "duplicate", err.Error())
}

func (s *PostgresTestSuite) TestCreate_SQLInjection() {
	err := s.repo.Create("", "")
	assert.Error(s.T(), err)
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
