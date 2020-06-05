package handlers

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			c.JSON(200, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		res, err := s.InsertUser(store.User{
			ID:       userID,
			Username: username,
			Password: string(hash),
		})
		if err != nil {
			c.JSON(500, gin.H{
				"success": res,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": res,
		})
	}
}

// Login logs in user
func Login(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")

		user, err := s.GetUserByUsername(username)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "Invalid credentials",
			})
			return
		}

		// Compare password hashes
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "Invalid credentials",
			})
			return
		}

		// TODO: supply token
		c.JSON(200, gin.H{
			"success": true,
		})
	}
}

// GetUser gets user by id
func GetUser(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Param("userID")
		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"user":    user,
		})
	}
}
