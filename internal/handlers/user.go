package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jasmaa/galt/internal/store"
)

// CreateAccount creates user account
func CreateAccount(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")

		// TODO: hash password here

		res := s.InsertUser(store.User{
			ID:       "54321", // TODO: generate uuid here
			Username: username,
			Password: password,
		})

		c.JSON(200, gin.H{
			"success": res,
		})
	}
}

// GetUser gets user by id
func GetUser(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Param("userID")
		user := s.GetUser(userID)

		c.JSON(200, gin.H{
			"user": user,
		})
	}
}
