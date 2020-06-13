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

func TestCreateAccountSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username=?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), "testuser", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create account
	data := url.Values{}
	data.Set("username", "testuser")
	data.Set("password", "testpassword")
	req, err := http.NewRequest("POST", "/api/v1/createAccount", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func TestCreateAccountFailNoCredentials(t *testing.T) {

	// Setup
	s := store.Store{}
	_ = s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	// Create account
	data := url.Values{}
	data.Set("username", "testuser")
	data.Set("password", "")
	req, err := http.NewRequest("POST", "/api/v1/createAccount", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username=?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Login
	data := url.Values{}
	data.Set("username", "testuser")
	data.Set("password", "testpassword")
	req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func TestLoginFailUnregisteredUser(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username=?").
		WithArgs("invaliduser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}))

	// Login
	data := url.Values{}
	data.Set("username", "invaliduser")
	data.Set("password", "testpassword")
	req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func TestLoginFailBadPassword(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username=?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Login
	data := url.Values{}
	data.Set("username", "testuser")
	data.Set("password", "invalidpassword")
	req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
