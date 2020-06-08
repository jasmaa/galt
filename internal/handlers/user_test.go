package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/jasmaa/galt/internal/store"
)

func TestGetUserSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Get user
	req, err := http.NewRequest("GET", "/api/v1/user/12345", nil)
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

func TestGetProfileSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Get user
	req, err := http.NewRequest("GET", "/api/v1/user", nil)
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

func TestUpdateProfileSuccess1(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))
	mock.ExpectExec("UPDATE users").
		WithArgs("12345", "leaf", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "I like to plant trees and eat eggplants.", "eggplant.png").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Get user
	data := url.Values{}
	data.Set("username", "leaf")
	data.Set("description", "I like to plant trees and eat eggplants.")
	data.Set("profileImgURL", "eggplant.png")
	req, err := http.NewRequest("PUT", "/api/v1/user", bytes.NewBufferString(data.Encode()))
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

func TestUpdateProfileSuccess2(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))
	mock.ExpectExec("UPDATE users").
		WithArgs("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", "eggplant.png").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Get user
	data := url.Values{}
	data.Set("profileImgURL", "eggplant.png")
	req, err := http.NewRequest("PUT", "/api/v1/user", bytes.NewBufferString(data.Encode()))
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

func TestDeleteProfileSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectExec("DELETE FROM users WHERE id=?").
		WithArgs("12345").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Get user
	req, err := http.NewRequest("DELETE", "/api/v1/user", nil)
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
