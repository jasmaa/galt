package handlers

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jasmaa/galt/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// CreateAccount creates user account
func CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)

		userID := uuid.New().String()
		username := c.PostForm("username")
		password := c.PostForm("password")

		// TODO: input validation

		if len(username) <= 0 || len(password) <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No username or password entered",
			})
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = s.InsertUser(store.User{
			ID:           userID,
			Username:     username,
			PasswordHash: string(hash),
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
func Login(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
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
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Create JWT token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": user.ID,
			"iat":    time.Now().Unix(),
			"exp":    time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		})
		tokenString, err := token.SignedString([]byte(hmacSecret))

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}
}
