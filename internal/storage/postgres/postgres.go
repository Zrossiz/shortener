package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostresStorage struct {
	conn *sql.DB
}

func New(conn *sql.DB) *PostresStorage {
	return &PostresStorage{conn: conn}
}

func Connect(DBURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to external DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping external DB: %w", err)
	}

	return db, nil
}

func (s *PostresStorage) Create(url string, hash string) error {
	query := `INSERT INTO urls (original, short) VALUES ($1, $2)`

	_, err := s.conn.Exec(query, url, hash)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("duplicate")
		}
		return fmt.Errorf("creating entry error: %w", err)
	}

	return nil
}

func (s *PostresStorage) Get(hash string) (string, error) {
	query := `SELECT original FROM urls WHERE short = $1`

	var original string
	err := s.conn.QueryRow(query, hash).Scan(&original)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("entry not found: %w", err)
		}
		return "", fmt.Errorf("error fetching entry: %w", err)
	}

	return original, nil
}
