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
	Poster          APIUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	Reshares        int       `form:"reshares" json:"reshares" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

// TODO: figure out response for auth vs not auth
// buildStatusResponse builds API status response
func buildStatusResponse(poster store.User, status store.Status, statusLikes int) APIStatus {
	return APIStatus{
		ID: status.ID,
		Poster: APIUser{
			ID:            poster.ID,
			Username:      poster.Username,
			Description:   poster.Description,
			ProfileImgURL: poster.ProfileImgURL,
		},
		Content:         status.Content,
		Likes:           statusLikes,
		Reshares:        -1,
		PostedTimestamp: status.PostedTimestamp,
		IsEdited:        status.IsEdited,
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

		statusLikes, err := s.GetStatusLikes(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes))
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

		c.JSON(http.StatusOK, buildStatusResponse(*user, status, 0))
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

		statusLikes, err := s.GetStatusLikes(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes))
	}
}

// LikeStatus likes status
func LikeStatus() gin.HandlerFunc {
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

		// Update status likes
		err = s.InsertStatusLikePair(userID, statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		statusLikes, err := s.GetStatusLikes(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes))
	}
}

// UnikeStatus unlikes status
func UnikeStatus() gin.HandlerFunc {
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

		// Update status likes
		err = s.DeleteStatusLikePair(userID, statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		statusLikes, err := s.GetStatusLikes(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes))
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

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status, 0))
	}
}
