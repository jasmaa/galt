package store

// Comment is a user comment
type Comment struct {
	ID              string
	UserID          string
	StatusID        string
	ParentCommentID string
}
