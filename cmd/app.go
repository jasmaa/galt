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

	v1.POST("/createAccount", handlers.CreateAccount())
	v1.POST("/login", handlers.Login(hmacSecret))

	v1.GET("/user/:userID", handlers.GetUser())
	v1.GET("/user", middleware.AuthUser(hmacSecret), handlers.GetProfile())
	v1.DELETE("/user", middleware.AuthUser(hmacSecret), handlers.DeleteProfile())

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
