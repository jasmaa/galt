package store

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestInsertUser(t *testing.T) {

	// Setup db
	s := Store{}
	s.Open()
	defer s.Close()

	user := User{
		ID:           "12345",
		Username:     "bowser",
		PasswordHash: "coolPasswordHash",
	}

	err := s.InsertUser(user)

	if err != nil {
		t.Errorf("error inserting user: %v", err)
	}
}
