package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jasmaa/galt/internal/store"
)

// CreateAccount creates user account
func CreateAccount(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := uuid.New().String()
		username := c.PostForm("username")
		password := c.PostForm("password")

		// TODO: input validation

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = s.InsertUser(store.User{
			ID:       userID,
			Username: username,
			Password: string(hash),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

// Login logs in user
func Login(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")

		user, err := s.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Compare password hashes
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Create session
		session := sessions.Default(c)
		session.Set("username", user.Username)
		session.Save()

		// TODO: supply token
		c.JSON(http.StatusOK, gin.H{})
	}
}

// GetUser gets user by id
func GetUser(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Param("userID")
		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}
