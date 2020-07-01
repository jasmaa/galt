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

		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		c.JSON(http.StatusOK, buildUserResponse(*authUser))
	}
}

// UpdateProfile updates user profile
func UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		username := c.PostForm("username")
		description := c.PostForm("description")
		profileImgURL := c.PostForm("profileImgURL")

		// Update user
		if len(username) > 0 {
			authUser.Username = username
		}
		if len(description) > 0 {
			authUser.Description = description
		}
		if len(profileImgURL) > 0 {
			authUser.ProfileImgURL = profileImgURL
		}

		err := s.UpdateUser(*authUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildUserResponse(*authUser))
	}
}

// DeleteProfile deletes user account
func DeleteProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		err := s.DeleteUserByID(authUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
