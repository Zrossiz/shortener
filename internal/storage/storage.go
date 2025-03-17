package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
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

func (s *Storage) Create(url string, hash string) (string, error) {
	query := `INSERT INTO urls (original, short) VALUES ($1, $2) RETURNING short`

	var short string
	err := s.conn.QueryRow(query, url, hash).Scan(&short)
	if err != nil {
		return "", fmt.Errorf("creating entry error: %w", err)
	}

	return short, nil
}

func (s *Storage) Get(hash string) (string, error) {
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
