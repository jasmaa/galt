package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jasmaa/galt/internal/store"
)

// APIUser represents a user in API response
type APIUser struct {
	ID            string `form:"id" json:"id" binding:"required"`
	Username      string `form:"username" json:"username" binding:"required"`
	Description   string `form:"description" json:"description" binding:"required"`
	ProfileImgURL string `form:"profileImgURL" json:"profileImgURL" binding:"required"`
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
			"user": APIUser{
				ID:            user.ID,
				Username:      user.Username,
				Description:   user.Description,
				ProfileImgURL: user.ProfileImgURL,
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
			"user": APIUser{
				ID:            user.ID,
				Username:      user.Username,
				Description:   user.Description,
				ProfileImgURL: user.ProfileImgURL,
			},
		})
	}
}

// DeleteProfile deletes user account
func DeleteProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		username := c.MustGet("username").(string)

		err := s.DeleteUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
