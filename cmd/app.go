package main

import (
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	cookieStore := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", cookieStore))

	// Handlers
	v1 := r.Group("/api/v1")
	{
		v1.POST("/createAccount", handlers.CreateAccount(s))
		v1.POST("/login", handlers.Login(s))
		v1.GET("/user/:userID", handlers.GetUser(s))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/panic", func(c *gin.Context) {
		panic("An unexpected error happen!")
	})
	r.GET("/testauth", middleware.AuthUser(), func(c *gin.Context) {
		session := sessions.Default(c)
		c.JSON(200, gin.H{
			"message": session.Get("username"),
		})
	})

	r.Run()
}
