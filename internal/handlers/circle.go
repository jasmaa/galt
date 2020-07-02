package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jasmaa/galt/internal/store"
)

// GetCircle gets circle by id
func GetCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		circleID := c.Param("circleID")
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		circle, err := s.GetCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if authUser.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to view this circle",
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(s, *circle))
	}
}

// CreateCircle creates new circle
func CreateCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		circleID := uuid.New().String()
		name := c.PostForm("name")
		description := c.PostForm("description")

		// Insert circle
		circle := store.Circle{
			ID:          circleID,
			UserID:      authUser.ID,
			Name:        name,
			Description: description,
		}
		err := s.InsertCircle(circle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(s, circle))
	}
}

// UpdateCircle updates circle
func UpdateCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		circleID := c.Param("circleID")
		name := c.PostForm("name")
		description := c.PostForm("description")

		circle, err := s.GetCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if authUser.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this circle",
			})
			return
		}

		// Update circle
		if len(name) > 0 {
			circle.Name = name
		}
		if len(description) > 0 {
			circle.Description = description
		}

		err = s.UpdateCircle(*circle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(s, *circle))
	}
}

// DeleteCircle deletes circle
func DeleteCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		circleID := c.Param("circleID")
		circle, err := s.GetCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if authUser.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to delete this circle",
			})
			return
		}

		// Delete status
		err = s.DeleteCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(s, *circle))
	}
}

// AddUserToCircle adds user to circle
func AddUserToCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		circleID := c.Param("circleID")
		userID := c.PostForm("userID")

		// Target user cannot be auth user
		if authUser.ID == userID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot add self to circle",
			})
			return
		}

		circle, err := s.GetCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if authUser.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this circle",
			})
			return
		}

		// Add target user to circle
		targetUser, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = s.InsertCircleUserPair(targetUser.ID, circle.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(s, *circle))
	}
}

// RemoveUserFromCircle removes user from circle
func RemoveUserFromCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		circleID := c.Param("circleID")
		userID := c.PostForm("userID")

		// Target user cannot be auth user
		if authUser.ID == userID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot add self to circle",
			})
			return
		}

		circle, err := s.GetCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if authUser.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this circle",
			})
			return
		}

		// Remove target user from circle
		targetUser, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = s.DeleteCircleUserPair(targetUser.ID, circle.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(s, *circle))
	}
}
