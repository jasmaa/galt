package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jasmaa/galt/internal/store"
)

// UserAPI represents a user in API response
type UserAPI struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
}

// GetUser gets user by id
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		userID := c.Param("userID")

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": UserAPI{
				ID:       user.ID,
				Username: user.Username,
			},
		})
	}
}

// GetProfile gets the user's profile
func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		username := c.MustGet("username").(string)

		user, err := s.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": UserAPI{
				ID:       user.ID,
				Username: user.Username,
			},
		})
	}
}
