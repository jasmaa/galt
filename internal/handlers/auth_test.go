package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

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

	v1.POST("/createAccount", CreateAccount())
	v1.POST("/login", Login(hmacSecret))

	v1.GET("/user/:userID", GetUser())
	v1.GET("/user", middleware.AuthUser(hmacSecret), GetProfile())
	v1.PUT("/user", middleware.AuthUser(hmacSecret), UpdateProfile())
	v1.DELETE("/user", middleware.AuthUser(hmacSecret), DeleteProfile())

	v1.GET("/status/:statusID", GetStatus())
	v1.POST("/status", middleware.AuthUser(hmacSecret), PostStatus())
	v1.PUT("/status/:statusID", middleware.AuthUser(hmacSecret), UpdateStatus())
	v1.DELETE("/status/:statusID", middleware.AuthUser(hmacSecret), DeleteStatus())

	return r
}

func TestCreateAccountAndLogin(t *testing.T) {

	// Setup
	s := store.Store{}
	s.Open()
	defer s.Close()
	r := setupRouter(s)

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

	// Login
	data = url.Values{}
	data.Set("username", "testuser")
	data.Set("password", "testpassword")
	req, err = http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}
