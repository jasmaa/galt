package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jasmaa/galt/internal/store"
)

// apiStatus represents a status in API response
type apiStatus struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          apiUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	IsLiked         bool      `form:"isLiked" json:"isLiked" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	Reshares        int       `form:"reshares" json:"reshares" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

// buildStatusResponse builds API status response
func buildStatusResponse(poster store.User, status store.Status, statusLikes int, isLiked bool, isReshared bool) apiStatus {
	return apiStatus{
		ID: status.ID,
		Poster: apiUser{
			ID:            poster.ID,
			Username:      poster.Username,
			Description:   poster.Description,
			ProfileImgURL: poster.ProfileImgURL,
		},
		Content:         status.Content,
		Likes:           statusLikes,
		IsLiked:         isLiked,
		Reshares:        -1, // filler value for now
		PostedTimestamp: status.PostedTimestamp,
		IsEdited:        status.IsEdited,
	}
}

// GetStatus gets status by id
func GetStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		statusID := c.Param("statusID")
		authUserID := c.MustGet("authUserID").(string)

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

		if len(authUserID) == 0 {
			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, false, false))
		} else {

			isLiked, err := s.GetIsStatusLiked(authUserID, statusID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, isLiked, false))
		}
	}
}

// PostStatus posts new status
func PostStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)

		statusID := uuid.New().String()
		authUserID := c.MustGet("authUserID").(string)
		content := c.PostForm("content")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		// TODO: Add content filtering here??

		// Insert status
		status := store.Status{
			ID:              statusID,
			UserID:          authUserID,
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

		user, err := s.GetUserByID(authUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(*user, status, 0, false, false))
	}
}

// UpdateStatus updates status
func UpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)
		statusID := c.Param("statusID")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		content := c.PostForm("content")

		status, err := s.GetStatusByID(statusID)
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

		if len(authUserID) == 0 {
			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, false, false))
		} else {

			isLiked, err := s.GetIsStatusLiked(authUserID, statusID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, isLiked, false))
		}
	}
}

// LikeStatus likes status
func LikeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)
		statusID := c.Param("statusID")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		status, err := s.GetStatusByID(statusID)
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

		// Update status likes
		err = s.InsertStatusLikePair(authUserID, statusID)
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

		if len(authUserID) == 0 {
			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, false, false))
		} else {

			isLiked, err := s.GetIsStatusLiked(authUserID, statusID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, isLiked, false))
		}
	}
}

// UnikeStatus unlikes status
func UnikeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)
		statusID := c.Param("statusID")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		status, err := s.GetStatusByID(statusID)
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

		// Update status likes
		err = s.DeleteStatusLikePair(authUserID, statusID)
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

		if len(authUserID) == 0 {
			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, false, false))
		} else {

			isLiked, err := s.GetIsStatusLiked(authUserID, statusID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, buildStatusResponse(*user, *status, statusLikes, isLiked, false))
		}
	}
}

// DeleteStatus updates status
func DeleteStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)
		statusID := c.Param("statusID")

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		status, err := s.GetStatusByID(statusID)
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

		c.JSON(http.StatusOK, buildStatusResponse(*user, *status, 0, false, false))
	}
}
