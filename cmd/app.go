package main

import (
	"log"
	"time"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jasmaa/galt/internal/handlers"
	"github.com/jasmaa/galt/internal/store"
)

func init() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	// TODO: setup db
	s := store.Store{}
	s.Open()

	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Router
	r := gin.New()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/user/:userID", handlers.GetUser(s))
		v1.POST("/user", handlers.CreateAccount(s))
	}

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
