package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
)

func main() {
	migrationsFolderPath := "./migrations"
	_ = godotenv.Load()
	DBURI := os.Getenv("DATABASE_DSN")

	conn, err := sql.Open("postgres", DBURI)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
	defer conn.Close()

	entries, err := os.ReadDir(migrationsFolderPath)
	if err != nil {
		panic(fmt.Errorf("failed to read migrations directory: %w", err))
	}

	for _, entry := range entries {
		sqlQuery, err := os.ReadFile(filepath.Join(migrationsFolderPath, entry.Name()))
		if err != nil {
			panic(fmt.Errorf("error reading SQL file %s: %w", entry.Name(), err))
		}

		_, err = conn.Exec(string(sqlQuery))
		if err != nil {
			panic(fmt.Errorf("error executing SQL from file %s: %w", entry.Name(), err))
		}

		fmt.Printf("Migration %s applied successfully\n", entry.Name())
	}
}
