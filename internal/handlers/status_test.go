package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/jasmaa/galt/internal/store"
)

func TestGetStatusSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM statuses WHERE id=?").
		WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "likes", "reshares", "posted_timestamp", "is_edited"}).
			AddRow("abcde", "12345", "I posted this status", 0, 0, time.Now(), false))
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Get status
	req, err := http.NewRequest("GET", "/api/v1/status/abcde", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostStatusSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectExec("INSERT INTO statuses").
		WithArgs(sqlmock.AnyArg(), "12345", "I got a bear today!", 0, 0, AnyTime{}, false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Post status
	data := url.Values{}
	data.Set("content", "I got a bear today!")
	req, err := http.NewRequest("POST", "/api/v1/status", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxMjM0NSIsImV4cCI6Ijk5OTk5OTk5OTk5OTk5OTkiLCJpYXQiOjE1MTYyMzkwMjJ9.M0It6-c6VJITiXp6WVgC-tNJdX5uB-YYqN2uJ75uzto")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateStatusSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM statuses WHERE id=?").
		WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "likes", "reshares", "posted_timestamp", "is_edited"}).
			AddRow("abcde", "12345", "I posted this status", 0, 0, time.Now(), false))
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))
	mock.ExpectExec("UPDATE statuses").
		WithArgs("abcde", "I got a bear today!", AnyTime{}, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Post status
	data := url.Values{}
	data.Set("content", "I got a bear today!")
	req, err := http.NewRequest("PUT", "/api/v1/status/abcde", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxMjM0NSIsImV4cCI6Ijk5OTk5OTk5OTk5OTk5OTkiLCJpYXQiOjE1MTYyMzkwMjJ9.M0It6-c6VJITiXp6WVgC-tNJdX5uB-YYqN2uJ75uzto")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateStatusFailUnauthorized(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM statuses WHERE id=?").
		WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "likes", "reshares", "posted_timestamp", "is_edited"}).
			AddRow("abcde", "67890", "this is my first post", 0, 0, time.Now(), false))
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Post status
	data := url.Values{}
	data.Set("content", "I got a bear today!")
	req, err := http.NewRequest("PUT", "/api/v1/status/abcde", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxMjM0NSIsImV4cCI6Ijk5OTk5OTk5OTk5OTk5OTkiLCJpYXQiOjE1MTYyMzkwMjJ9.M0It6-c6VJITiXp6WVgC-tNJdX5uB-YYqN2uJ75uzto")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteStatusSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM statuses WHERE id=?").
		WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "likes", "reshares", "posted_timestamp", "is_edited"}).
			AddRow("abcde", "12345", "I posted this status", 0, 0, time.Now(), false))
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))
	mock.ExpectExec("DELETE FROM statuses").
		WithArgs("abcde").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Delete status
	req, err := http.NewRequest("DELETE", "/api/v1/status/abcde", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxMjM0NSIsImV4cCI6Ijk5OTk5OTk5OTk5OTk5OTkiLCJpYXQiOjE1MTYyMzkwMjJ9.M0It6-c6VJITiXp6WVgC-tNJdX5uB-YYqN2uJ75uzto")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteStatusFailUnauthorized(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM statuses WHERE id=?").
		WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content", "likes", "reshares", "posted_timestamp", "is_edited"}).
			AddRow("abcde", "67890", "I posted this status", 0, 0, time.Now(), false))
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Delete status
	req, err := http.NewRequest("DELETE", "/api/v1/status/abcde", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiIxMjM0NSIsImV4cCI6Ijk5OTk5OTk5OTk5OTk5OTkiLCJpYXQiOjE1MTYyMzkwMjJ9.M0It6-c6VJITiXp6WVgC-tNJdX5uB-YYqN2uJ75uzto")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
