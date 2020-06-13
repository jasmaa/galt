package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jasmaa/galt/internal/store"
)

// APIStatus represents a status in API response
type APIStatus struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	Reshares        int       `form:"reshares" json:"reshares" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

func buildStatusResponse(user store.User, status store.Status) gin.H {
	return gin.H{
		"user": APIUser{
			ID:            user.ID,
			Username:      user.Username,
			Description:   user.Description,
			ProfileImgURL: user.ProfileImgURL,
		},
		"status": APIStatus{
			ID:              status.ID,
			Content:         status.Content,
			Likes:           status.Likes,
			Reshares:        status.Reshares,
			PostedTimestamp: status.PostedTimestamp,
			IsEdited:        status.IsEdited,
		},
	}
}

// GetStatus gets status by id
func GetStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		statusID := c.Param("statusID")

		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.GetUserByID(status.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status))
	}
}

// PostStatus posts new status
func PostStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)

		statusID := uuid.New().String()
		userID := c.MustGet("userID").(string)
		content := c.PostForm("content")

		// TODO: Add content filtering here??

		// Insert status
		status := store.Status{
			ID:              statusID,
			UserID:          userID,
			Content:         content,
			Likes:           0,
			Reshares:        0,
			PostedTimestamp: time.Now(),
			IsEdited:        false,
		}
		err := s.InsertStatus(status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, status))
	}
}

// UpdateStatus updates status
func UpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		userID := c.MustGet("userID").(string)
		statusID := c.Param("statusID")

		content := c.PostForm("content")

		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user is poster
		if status.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to edit this status",
			})
			return
		}

		// TODO: validate content

		// Update status
		status.Content = content
		status.PostedTimestamp = time.Now()
		status.IsEdited = true

		err = s.UpdateStatus(*status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status))
	}
}

// DeleteStatus updates status
func DeleteStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		userID := c.MustGet("userID").(string)
		statusID := c.Param("statusID")

		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user is poster
		if status.UserID != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You do not have permission to delete this status",
			})
			return
		}

		// Delete status
		err = s.DeleteStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status))
	}
}
