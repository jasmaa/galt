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

// buildUserResponse builds API user response
func buildUserResponse(user store.User) APIUser {
	return APIUser{
		ID:            user.ID,
		Username:      user.Username,
		Description:   user.Description,
		ProfileImgURL: user.ProfileImgURL,
	}
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

		c.JSON(http.StatusOK, buildUserResponse(*user))
	}
}

// GetProfile gets the user's profile
func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		userID := c.MustGet("userID").(string)

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildUserResponse(*user))
	}
}

// UpdateProfile updates user profile
func UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		userID := c.MustGet("userID").(string)

		username := c.PostForm("username")
		description := c.PostForm("description")
		profileImgURL := c.PostForm("profileImgURL")

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Update user
		if len(username) > 0 {
			user.Username = username
		}
		if len(description) > 0 {
			user.Description = description
		}
		if len(profileImgURL) > 0 {
			user.ProfileImgURL = profileImgURL
		}

		err = s.UpdateUser(*user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildUserResponse(*user))
	}
}

// DeleteProfile deletes user account
func DeleteProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		userID := c.MustGet("userID").(string)

		err := s.DeleteUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
