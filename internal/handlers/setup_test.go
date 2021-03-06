package handlers_test

import (
	"database/sql/driver"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jasmaa/galt/internal/handlers"
	"github.com/jasmaa/galt/internal/middleware"
	"github.com/jasmaa/galt/internal/store"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// setupRouter sets up router for test
func setupRouter(s store.Store) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	hmacSecret := "secret_key"

	// Router
	r := gin.New()

	// Handlers
	v1 := r.Group("/api/v1")

	v1.Use(func(c *gin.Context) {
		c.Set("store", s)
		c.Next()
	})

	v1.Use(middleware.AuthUser(hmacSecret))

	v1.POST("/createAccount", handlers.CreateAccount())
	v1.POST("/login", handlers.Login(hmacSecret))

	v1.GET("/user/:userID", handlers.GetUser())
	v1.GET("/user", handlers.GetProfile())
	v1.PUT("/user", handlers.UpdateProfile())
	v1.DELETE("/user", handlers.DeleteProfile())

	v1.GET("/status/:statusID", handlers.GetStatus())
	v1.POST("/status", handlers.PostStatus())
	v1.PUT("/status/:statusID", handlers.UpdateStatus())
	v1.DELETE("/status/:statusID", handlers.DeleteStatus())
	v1.POST("/status/:statusID/like", handlers.LikeStatus())
	v1.POST("/status/:statusID/unlike", handlers.UnikeStatus())
	v1.GET("/status/:statusID/comments", handlers.GetComments())
	v1.POST("/status/:statusID/comment", handlers.PostComment())

	return r
}
