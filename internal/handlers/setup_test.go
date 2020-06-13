package handlers_test

import (
	"github.com/gin-gonic/gin"

	"github.com/jasmaa/galt/internal/handlers"
	"github.com/jasmaa/galt/internal/middleware"
	"github.com/jasmaa/galt/internal/store"
)

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
