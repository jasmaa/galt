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
		authUserID := c.MustGet("authUserID").(string)

		circle, err := s.GetCircleByID(circleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
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

		// Check if user owns circle
		if user.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this status",
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(*circle))
	}
}

// CreateCircle creates new circle
func CreateCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)

		circleID := uuid.New().String()
		authUserID := c.MustGet("authUserID").(string)
		name := c.PostForm("name")
		description := c.PostForm("description")

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

		// Insert circle
		circle := store.Circle{
			ID:          circleID,
			UserID:      user.ID,
			Name:        name,
			Description: description,
		}
		err = s.InsertCircle(circle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCircleResponse(circle))
	}
}

// UpdateCircle updates circle
func UpdateCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)
		circleID := c.Param("circleID")
		name := c.PostForm("name")
		description := c.PostForm("description")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
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

		user, err := s.GetUserByID(authUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if user.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this status",
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

		c.JSON(http.StatusOK, buildCircleResponse(*circle))
	}
}

// DeleteCircle deletes circle
func DeleteCircle() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)
		circleID := c.Param("circleID")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
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

		user, err := s.GetUserByID(authUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user owns circle
		if user.ID != circle.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this status",
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

		c.JSON(http.StatusOK, buildCircleResponse(*circle))
	}
}
