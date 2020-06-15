package store

import (
	"errors"
	"time"
)

// Status is a status post
type Status struct {
	ID              string
	UserID          string
	Content         string
	PostedTimestamp time.Time
	IsEdited        bool
}

// GetStatusByID gets status by id
func (s *Store) GetStatusByID(statusID string) (*Status, error) {

	row := s.db.QueryRow(
		"SELECT id, user_id, content, posted_timestamp, is_edited FROM statuses WHERE id=$1",
		statusID,
	)

	status := Status{}
	if err := row.Scan(&status.ID, &status.UserID, &status.Content, &status.PostedTimestamp, &status.IsEdited); err != nil {
		return nil, errors.New("No status found")
	}

	return &status, nil
}

// GetStatusFeed gets statuses in feed
func (s *Store) GetStatusFeed(userID string) ([]Status, error) {

	rows, err := s.db.Query(
		`SELECT id, user_id, content, posted_timestamp, is_edited
		FROM statuses WHERE user_id IN (
			SELECT circle_user_pairs.user_id
			FROM circle_user_pairs JOIN circles ON circle_user_pairs.circle_id=circles.id
			WHERE circles.user_id=$1
		) ORDER BY posted_timestamp DESC`,
		userID,
	)
	defer rows.Close()
	if err != nil {
		return nil, errors.New("Error retrieving feed")
	}

	statuses := make([]Status, 0)

	for rows.Next() {
		status := Status{}
		rows.Scan(&status.ID, &status.UserID, &status.Content, &status.PostedTimestamp, &status.IsEdited)
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// InsertStatus inserts status
func (s *Store) InsertStatus(status Status) error {

	_, err := s.db.Exec(
		"INSERT INTO statuses (id, user_id, content, posted_timestamp, is_edited) VALUES ($1, $2, $3, $4, $5)",
		status.ID, status.UserID, status.Content, status.PostedTimestamp, status.IsEdited,
	)
	if err != nil {
		return errors.New("Error creating status")
	}

	return nil
}

// UpdateStatus udpates status
func (s *Store) UpdateStatus(status Status) error {

	_, err := s.db.Exec(
		"UPDATE statuses SET content=$2, posted_timestamp=$3, is_edited=$4 WHERE id=$1",
		status.ID, status.Content, status.PostedTimestamp, status.IsEdited,
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

// GetStatusLikes gets number of likes on a status
func (s *Store) GetStatusLikes(statusID string) (int, error) {

	row := s.db.QueryRow("SELECT COUNT(status_id) FROM status_like_pairs WHERE status_id=$1",
		statusID,
	)
	var count int
	if err := row.Scan(&count); err != nil {
		return -1, errors.New("Error with database")
	}

	return count, nil
}

// GetIsUserLikedStatus checks if status was liked by user
func (s *Store) GetIsUserLikedStatus(userID string, statusID string) (bool, error) {

	row := s.db.QueryRow("SELECT COUNT(status_id) FROM status_like_pairs WHERE user_id=$1 AND status_id=$2",
		userID, statusID,
	)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, errors.New("Error with database")
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

// InsertStatusLikePair inserts (userID, statusID) pair for liking a status post
func (s *Store) InsertStatusLikePair(userID string, statusID string) error {

	// Check for duplicate
	isLiked, err := s.GetIsUserLikedStatus(userID, statusID)
	if err != nil {
		return err
	}

	if isLiked {
		return nil
	}

	_, err = s.db.Exec(
		"INSERT INTO status_like_pairs (user_id, status_id) VALUES ($1, $2)",
		userID, statusID,
	)
	if err != nil {
		return errors.New("Error inserting status like pair")
	}

	return nil
}

// DeleteStatusLikePair deletes (userID, statusID) pair for liking a status post
func (s *Store) DeleteStatusLikePair(userID string, statusID string) error {

	_, err := s.db.Exec(
		"DELETE FROM status_like_pairs WHERE user_id=$1 AND status_id=$2",
		userID, statusID,
	)
	if err != nil {
		return errors.New("Error deleting status like pair")
	}

	return nil
}
