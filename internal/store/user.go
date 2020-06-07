package store

import "errors"

// User is a site user
type User struct {
	ID           string
	Username     string
	PasswordHash string
}

// InsertUser inserts user into db
func (s *Store) InsertUser(user User) error {

	// Check for duplicate username
	row := s.db.QueryRow("SELECT COUNT(username) FROM users WHERE username=$1", user.Username)
	var count int
	if err := row.Scan(&count); err != nil {
		return errors.New("Error with database")
	}
	if count > 0 {
		return errors.New("User already exists")
	}

	// Insert user
	_, err := s.db.Exec(
		"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		user.ID, user.Username, user.PasswordHash,
	)
	if err != nil {
		return errors.New("Error creating user")
	}

	return nil
}

// GetUserByID gets user by id from db
func (s *Store) GetUserByID(id string) (*User, error) {

	row := s.db.QueryRow("SELECT id, username, password FROM users WHERE id=$1", id)

	user := User{}
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		return nil, errors.New("No user found")
	}

	return &user, nil
}

// GetUserByUsername gets user by username from db
func (s *Store) GetUserByUsername(username string) (*User, error) {

	row := s.db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username)

	user := User{}
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		return nil, errors.New("No user found")
	}

	return &user, nil
}

// DeleteUserByUsername deletes user by username
func (s *Store) DeleteUserByUsername(username string) error {

	_, err := s.db.Exec(
		"DELETE FROM users WHERE username=$1",
		username,
	)
	if err != nil {
		return errors.New("Error deleting user")
	}

	return nil
}
