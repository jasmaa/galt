package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jasmaa/galt/internal/store"
)

// APIStatus represents a status in API response
type APIStatus struct {
	ID      string `form:"id" json:"id" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
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

		c.JSON(http.StatusOK, gin.H{
			"user": APIUser{
				ID:            user.ID,
				Username:      user.Username,
				Description:   user.Description,
				ProfileImgURL: user.ProfileImgURL,
			},
			"status": APIStatus{
				ID:      status.ID,
				Content: status.Content,
			},
		})
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
		err := s.InsertStatus(store.Status{
			ID:      statusID,
			UserID:  userID,
			Content: content,
		})
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

		c.JSON(http.StatusOK, gin.H{
			"user": APIUser{
				ID:            user.ID,
				Username:      user.Username,
				Description:   user.Description,
				ProfileImgURL: user.ProfileImgURL,
			},
			"status": APIStatus{
				ID:      statusID,
				Content: content,
			},
		})
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

		// Update status
		if len(content) > 0 {
			status.Content = content
		}

		err = s.UpdateStatus(*status)
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
			"status": APIStatus{
				ID:      status.ID,
				Content: status.Content,
			},
		})
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

		c.JSON(http.StatusOK, gin.H{
			"user": APIUser{
				ID:            user.ID,
				Username:      user.Username,
				Description:   user.Description,
				ProfileImgURL: user.ProfileImgURL,
			},
			"status": APIStatus{
				ID:      status.ID,
				Content: status.Content,
			},
		})
	}
}
