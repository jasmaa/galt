package store

import (
	"errors"
)

// Status is a status post
type Status struct {
	ID      string
	UserID  string
	Content string
}

// GetStatusByID gets status by id
func (s *Store) GetStatusByID(statusID string) (*Status, error) {

	row := s.db.QueryRow("SELECT id, user_id, content FROM statuses WHERE id=$1", statusID)

	status := Status{}
	if err := row.Scan(&status.ID, &status.UserID, &status.Content); err != nil {
		return nil, errors.New("No status found")
	}

	return &status, nil
}

// InsertStatus inserts status
func (s *Store) InsertStatus(status Status) error {

	_, err := s.db.Exec(
		"INSERT INTO statuses (id, user_id, content) VALUES ($1, $2, $3)",
		status.ID, status.UserID, status.Content,
	)
	if err != nil {
		return errors.New("Error creating status")
	}

	return nil
}

// UpdateStatus udpates status
func (s *Store) UpdateStatus(status Status) error {

	_, err := s.db.Exec(
		"UPDATE statuses SET content=$2 WHERE id=$1",
		status.ID, status.Content,
	)
	if err != nil {
		return errors.New("Error updating status")
	}

	return nil
}

// DeleteStatusByID deletes status by statusID
func (s *Store) DeleteStatusByID(statusID string) error {

	_, err := s.db.Exec("DELETE FROM statuses WHERE id=$1", statusID)
	if err != nil {
		return errors.New("Error deleting status")
	}

	return nil
}
