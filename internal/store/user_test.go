package store_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"

	"github.com/jasmaa/galt/internal/store"
)

func TestInsertUser(t *testing.T) {

	// Setup db
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs("bowser").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO users").
		WithArgs("12345", "bowser", "coolPasswordHash").
		WillReturnResult(sqlmock.NewResult(0, 0))

	user := store.User{
		ID:           "12345",
		Username:     "bowser",
		PasswordHash: "coolPasswordHash",
	}

	err := s.InsertUser(user)

	if err != nil {
		t.Errorf("error inserting user: %v", err)
	}
}
