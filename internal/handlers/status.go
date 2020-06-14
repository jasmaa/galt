package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jasmaa/galt/internal/store"
)

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

			isLiked, err := s.GetIsUserLikedStatus(authUserID, statusID)
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

// GetStatusFeed gets users feed
func GetStatusFeed() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUserID := c.MustGet("authUserID").(string)

		if len(authUserID) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No token provided",
			})
			return
		}

		_, err := s.GetUserByID(authUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		statuses, err := s.GetStatusFeed(authUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		apiStatuses := make(map[string]interface{})
		for k, v := range statuses {

			user, err := s.GetUserByID(v.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			statusLikes, err := s.GetStatusLikes(v.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			isLiked, err := s.GetIsUserLikedStatus(authUserID, v.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			apiStatuses[k] = buildStatusResponse(*user, v, statusLikes, isLiked, false)
		}

		c.JSON(http.StatusOK, apiStatuses)
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

			isLiked, err := s.GetIsUserLikedStatus(authUserID, statusID)
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

			isLiked, err := s.GetIsUserLikedStatus(authUserID, statusID)
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

			isLiked, err := s.GetIsUserLikedStatus(authUserID, statusID)
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

// GetComments gets comments on a status
func GetComments() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		statusID := c.Param("statusID")

		comments, err := s.GetCommentsFromStatus(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO: do join instead??? or make client find user info instead??
		apiComments := make(map[string]interface{})
		for k, v := range comments {

			user, err := s.GetUserByID(v.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			commentLikes, err := s.GetCommentLikes(v.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			apiComments[k] = buildCommentResponse(*user, v, commentLikes)
		}

		c.JSON(http.StatusOK, apiComments)
	}
}

// PostComment posts a comment on a status
func PostComment() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		commentID := uuid.New().String()
		authUserID := c.MustGet("authUserID").(string)
		content := c.PostForm("content")
		statusID := c.Param("statusID")

		comment := store.Comment{
			ID:              commentID,
			UserID:          authUserID,
			StatusID:        statusID,
			ParentCommentID: sql.NullString{},
			Content:         content,
			PostedTimestamp: time.Now(),
			IsEdited:        false,
		}
		err := s.InsertComment(comment)
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

		c.JSON(http.StatusOK, buildCommentResponse(*user, comment, 0))
	}
}
