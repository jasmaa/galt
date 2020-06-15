package store

import (
	"database/sql"
	"errors"
	"time"
)

// Comment is a user comment
type Comment struct {
	ID              string
	UserID          string
	StatusID        string
	ParentCommentID sql.NullString
	Content         string
	PostedTimestamp time.Time
	IsEdited        bool
}

// GetCommentByID gets comment by id
func (s *Store) GetCommentByID(commentID string) (*Comment, error) {

	row := s.db.QueryRow(
		`SELECT id, user_id, status_id, parent_comment_id, content, posted_timestamp, is_edited FROM comments WHERE id=$1`,
		commentID,
	)

	comment := Comment{}
	if err := row.Scan(&comment.ID, &comment.UserID, &comment.StatusID, &comment.ParentCommentID, &comment.Content, &comment.PostedTimestamp, &comment.IsEdited); err != nil {
		return nil, errors.New("No comment found")
	}

	return &comment, nil
}

// GetCommentsFromStatus gets comments under a status
func (s *Store) GetCommentsFromStatus(statusID string) ([]Comment, error) {

	rows, err := s.db.Query(
		`SELECT id, user_id, status_id, parent_comment_id, content, posted_timestamp, is_edited
		FROM comments WHERE status_id=$1
		ORDER BY posted_timestamp DESC`,
		statusID,
	)
	defer rows.Close()
	if err != nil {
		return nil, errors.New("Error retrieving comments")
	}

	comments := make([]Comment, 0)

	for rows.Next() {
		comment := Comment{}
		rows.Scan(&comment.ID, &comment.UserID, &comment.StatusID, &comment.ParentCommentID, &comment.Content, &comment.PostedTimestamp, &comment.IsEdited)
		comments = append(comments, comment)
	}

	return comments, nil
}

// InsertComment inserts comment
func (s *Store) InsertComment(comment Comment) error {

	_, err := s.db.Exec(
		"INSERT INTO comments (id, user_id, status_id, parent_comment_id, content, posted_timestamp, is_edited) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		comment.ID, comment.UserID, comment.StatusID, comment.ParentCommentID, comment.Content, comment.PostedTimestamp, comment.IsEdited,
	)
	if err != nil {
		return errors.New("Error creating comment")
	}

	return nil
}

// UpdateComment updates comment
func (s *Store) UpdateComment(comment Comment) error {

	_, err := s.db.Exec(
		"UPDATE comments SET content=$2, posted_timestamp=$3, is_edited=$4 WHERE id=$1",
		comment.ID, comment.Content, comment.PostedTimestamp, comment.IsEdited,
	)
	if err != nil {
		return errors.New("Error updating comment")
	}

	return nil
}

// DeleteCommentByID deletes comment by commentID
func (s *Store) DeleteCommentByID(commentID string) error {

	_, err := s.db.Exec("DELETE FROM comments WHERE id=$1", commentID)
	if err != nil {
		return errors.New("Error deleting comment")
	}

	return nil
}

// GetCommentLikes gets number of likes on a comment
func (s *Store) GetCommentLikes(commentID string) (int, error) {

	row := s.db.QueryRow("SELECT COUNT(comment_id) FROM comment_like_pairs WHERE comment_id=$1",
		commentID,
	)
	var count int
	if err := row.Scan(&count); err != nil {
		return -1, errors.New("Error with database")
	}

	return count, nil
}

// GetIsUserLikedComment checks if comment was liked by user
func (s *Store) GetIsUserLikedComment(userID string, commentID string) (bool, error) {

	row := s.db.QueryRow("SELECT COUNT(comment_id) FROM comment_like_pairs WHERE user_id=$1 AND comment_id=$2",
		userID, commentID,
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

// InsertCommentLikePair inserts (userID, commentID) pair for liking a comment
func (s *Store) InsertCommentLikePair(userID string, commentID string) error {

	// Check for duplicate
	isLiked, err := s.GetIsUserLikedComment(userID, commentID)
	if err != nil {
		return err
	}

	if isLiked {
		return nil
	}

	_, err = s.db.Exec(
		"INSERT INTO comment_like_pairs (user_id, comment_id) VALUES ($1, $2)",
		userID, commentID,
	)
	if err != nil {
		return errors.New("Error inserting comment like pair")
	}

	return nil
}

// DeleteCommentLikePair deletes (userID, commentID) pair for liking a comment
func (s *Store) DeleteCommentLikePair(userID string, commentID string) error {

	_, err := s.db.Exec(
		"DELETE FROM comment_like_pairs WHERE user_id=$1 AND comment_id=$2",
		userID, commentID,
	)
	if err != nil {
		return errors.New("Error deleting comment like pair")
	}

	return nil
}
