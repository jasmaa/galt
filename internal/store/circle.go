package store

import (
	"errors"
)

// Circle is a social circle
type Circle struct {
	ID          string
	UserID      string
	Name        string
	Description string
}

// GetCircleByID gets circle by id
func (s *Store) GetCircleByID(circleID string) (*Circle, error) {

	row := s.db.QueryRow(
		`SELECT id, user_id, name, description FROM circles WHERE id=$1`,
		circleID,
	)

	circle := Circle{}
	if err := row.Scan(&circle.ID, &circle.UserID, &circle.Name, &circle.Description); err != nil {
		return nil, errors.New("No circle found")
	}

	return &circle, nil
}

// InsertCircle inserts social circle
func (s *Store) InsertCircle(circle Circle) error {

	_, err := s.db.Exec(
		"INSERT INTO circles (id, user_id, name, description) VALUES ($1, $2, $3, $4)",
		circle.ID, circle.UserID, circle.Name, circle.Description,
	)
	if err != nil {
		return errors.New("Error creating circle")
	}

	return nil
}

// UpdateCircle updates circle
func (s *Store) UpdateCircle(circle Circle) error {

	_, err := s.db.Exec(
		"UPDATE circles SET name=$2, description=$3 WHERE id=$1",
		circle.ID, circle.Name, circle.Description,
	)
	if err != nil {
		return errors.New("Error updating circle")
	}

	return nil
}

// DeleteCircleByID deletes circle by circleID
func (s *Store) DeleteCircleByID(circleID string) error {

	_, err := s.db.Exec("DELETE FROM circles WHERE id=$1", circleID)
	if err != nil {
		return errors.New("Error deleting circle")
	}

	return nil
}

// GetIsInCircle checks if comment was liked by user
func (s *Store) GetIsInCircle(userID string, circleID string) (bool, error) {

	row := s.db.QueryRow("SELECT COUNT(circleID) FROM user_circle_pairs WHERE user_id=$1 AND circle_id=$2",
		userID, circleID,
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

// InsertCircleUserPair inserts (userID, circleID) pair for adding user to social circle
func (s *Store) InsertCircleUserPair(userID string, circleID string) error {

	// Check for duplicate
	isInCircle, err := s.GetIsInCircle(userID, circleID)
	if err != nil {
		return err
	}

	if isInCircle {
		return nil
	}

	_, err = s.db.Exec(
		"INSERT INTO user_circle_pairs (user_id, circle_id) VALUES ($1, $2)",
		userID, circleID,
	)
	if err != nil {
		return errors.New("Error inserting user circle pair")
	}

	return nil
}

// DeleteUserCirclePair deletes (userID, circleID) pair
func (s *Store) DeleteUserCirclePair(userID string, circleID string) error {

	_, err := s.db.Exec(
		"DELETE FROM user_circle_pairs WHERE user_id=$1 AND circle_id=$2",
		userID, circleID,
	)
	if err != nil {
		return errors.New("Error deleting user circle pair")
	}

	return nil
}
