package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/jasmaa/galt/internal/handlers"
	"github.com/jasmaa/galt/internal/middleware"
	"github.com/jasmaa/galt/internal/store"
)

func setupRouter(s store.Store) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	hmacSecret := "some_secret"

	// Router
	r := gin.New()

	// Handlers
	v1 := r.Group("/api/v1")

	v1.Use(func(c *gin.Context) {
		c.Set("store", s)
		c.Next()
	})

	v1.POST("/createAccount", handlers.CreateAccount())
	v1.POST("/login", handlers.Login(hmacSecret))

	v1.GET("/user/:userID", handlers.GetUser())
	v1.GET("/user", middleware.AuthUser(hmacSecret), handlers.GetProfile())
	v1.PUT("/user", middleware.AuthUser(hmacSecret), handlers.UpdateProfile())
	v1.DELETE("/user", middleware.AuthUser(hmacSecret), handlers.DeleteProfile())

	v1.GET("/status/:statusID", handlers.GetStatus())
	v1.POST("/status", middleware.AuthUser(hmacSecret), handlers.PostStatus())
	v1.PUT("/status/:statusID", middleware.AuthUser(hmacSecret), handlers.UpdateStatus())
	v1.DELETE("/status/:statusID", middleware.AuthUser(hmacSecret), handlers.DeleteStatus())

	return r
}

func TestCreateAccount(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), "testuser", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0))

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

func TestLoginSuccess(t *testing.T) {

	// Setup
	s := store.Store{}
	mock := s.OpenMock()
	defer s.Close()
	r := setupRouter(s)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "description", "profile_img_url"}).
			AddRow("12345", "testuser", "$2b$10$KpZAZIPai8SyT7k8zT582ec5Va9.KrnoMc9D5UnGkDRdVvTp263/q", "", ""))

	// Create account
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
