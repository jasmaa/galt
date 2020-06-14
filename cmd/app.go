package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jasmaa/galt/internal/handlers"
	"github.com/jasmaa/galt/internal/middleware"
	"github.com/jasmaa/galt/internal/store"
)

func init() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	hmacSecret := os.Getenv("HMAC_SECRET")

	// Setup db
	s := store.Store{}
	s.Open()
	defer s.Close()

	// Router
	r := gin.New()

	// Middleware
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

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

	/*
		v1.GET("/comment/:commentID", handlers.GetStatus())
		v1.PUT("/comment/:commentID", middleware.AuthUser(hmacSecret), handlers.UpdateStatus())
		v1.DELETE("/comment/:commentID", middleware.AuthUser(hmacSecret), handlers.DeleteStatus())
		v1.POST("/comment/:commentID/like", middleware.AuthUser(hmacSecret), handlers.LikeStatus())
		v1.POST("/comment/:commentID/unlike", middleware.AuthUser(hmacSecret), handlers.UnikeStatus())
	*/

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/panic", func(c *gin.Context) {
		panic("An unexpected error happen!")
	})

	r.Run()
}
