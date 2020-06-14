package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jasmaa/galt/internal/store"
)

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
		authUserID := c.MustGet("authUserID").(string)

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		user, err := s.GetUserByID(authUserID)
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
		authUserID := c.MustGet("authUserID").(string)

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		username := c.PostForm("username")
		description := c.PostForm("description")
		profileImgURL := c.PostForm("profileImgURL")

		user, err := s.GetUserByID(authUserID)
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
		authUserID := c.MustGet("authUserID").(string)

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		err := s.DeleteUserByID(authUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
