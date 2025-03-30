package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type PosgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PosgresRepo {
	return &PosgresRepo{db: conn}
}

func Connect(uri string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %v", err)
	}

	return conn, nil
}

func (p *PosgresRepo) Create(url string, hash string) error {
	query := `INSERT INTO urls (original, short) VALUES ($1, $2)`

	_, err := p.db.Exec(query, url, hash)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("duplicate")
		}
		return fmt.Errorf("creating entry error: %w", err)
	}

	return nil
}
