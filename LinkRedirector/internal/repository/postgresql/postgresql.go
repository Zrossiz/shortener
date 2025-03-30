package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

func (s *PosgresRepo) Get(hash string) (string, error) {
	query := `SELECT original FROM urls WHERE short = $1`

	var original string
	err := s.db.QueryRow(query, hash).Scan(&original)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("entry not found: %w", err)
		}
		return "", fmt.Errorf("error fetching entry: %w", err)
	}

	return original, nil
}
