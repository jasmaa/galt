package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Store wraps database connection
type Store struct {
	db *sql.DB
}

// Open opens database connection
func (s *Store) Open() {

	// Connect to db
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgDB := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPassword, pgHost, pgPort, pgDB,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	s.db = db
}

// Close closes database connection
func (s *Store) Close() {
	s.db.Close()
}
