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

func (s *Store) InsertUser(user User) bool {

	_, err := s.db.Exec(
		"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		user.ID, user.Username, user.Password,
	)

	if err != nil {
		return false
	}

	return true
}

// GetUser gets user with id from db
func (s *Store) GetUser(id string) *User {

	row := s.db.QueryRow("SELECT id, username, password FROM users WHERE id=$1", id)

	user := User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil
	}

	return &user
}
