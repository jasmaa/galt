package handlers

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
func Login(s store.Store, hmacSecret string) gin.HandlerFunc {
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

		// Create JWT token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		})
		tokenString, err := token.SignedString([]byte(hmacSecret))

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
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
