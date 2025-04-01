package clickhouse

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Zrossiz/LogConsumer/consumer/internal/domain"
)

type ClickhouseDB struct {
	conn *sql.DB
}

func ClickhouseConnect(uri string) (*sql.DB, error) {
	conn, err := sql.Open("clickhouse", uri)
	if err != nil {
		return nil, fmt.Errorf("error connect to db %v", err)
	}

	return conn, nil
}

func NewClickhouse(conn *sql.DB) *ClickhouseDB {
	return &ClickhouseDB{
		conn: conn,
	}
}

func (c *ClickhouseDB) Create(data domain.RegisterRedirectEvent) error {
	query := `INSERT INTO redirects (original, short, user_ip, os, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := c.conn.Exec(query, data.Original, data.Short, data.UserIP, data.Os, time.Now())
	if err != nil {
		return fmt.Errorf("insert redirect error: %v", err)
	}

	return nil
}

func (c *ClickhouseDB) Get() ([]domain.RedirectEventDAO, error) {
	query := `SELECT id, original, short, user_ip, os FROM redirects`

	rows, err := c.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error get redirects from db: %v", err)
	}
	defer rows.Close()

	var data []domain.RedirectEventDAO
	for rows.Next() {
		var redirect domain.RedirectEventDAO
		err := rows.Scan(
			&redirect.ID,
			&redirect.Original,
			&redirect.Short,
			&redirect.UserIP,
			&redirect.Os,
			&redirect.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scan row: %v", err)
		}

		data = append(data, redirect)
	}

	return data, nil
}
