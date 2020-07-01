package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
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

		authUser, _ := c.MustGet("authUser").(*store.User)

		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		poster, err := s.GetUserByID(status.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(s, *poster, *status, authUser))
	}
}

// GetStatusFeed gets users feed
func GetStatusFeed() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			offset = 0
		} else if offset < 0 {
			offset = 0
		}

		statuses, err := s.GetStatusFeed(authUser.ID, 30, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		apiStatuses := make([]interface{}, len(statuses))
		for i, status := range statuses {

			poster, err := s.GetUserByID(status.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			apiStatuses[i] = buildStatusResponse(s, *poster, status, authUser)
		}

		c.JSON(http.StatusOK, gin.H{
			"statuses": apiStatuses,
			"offset":   offset + 30,
		})
	}
}

// PostStatus posts new status
func PostStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		statusID := uuid.New().String()
		content := c.PostForm("content")

		// TODO: Add content filtering here??

		// Insert status
		status := store.Status{
			ID:              statusID,
			UserID:          authUser.ID,
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

		c.JSON(http.StatusOK, buildStatusResponse(s, *authUser, status, authUser))
	}
}

// UpdateStatus updates status
func UpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		statusID := c.Param("statusID")
		content := c.PostForm("content")

		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user is poster
		if status.UserID != authUser.ID {
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

		c.JSON(http.StatusOK, buildStatusResponse(s, *authUser, *status, authUser))
	}
}

// LikeStatus likes status
func LikeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		statusID := c.Param("statusID")
		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		poster, err := s.GetUserByID(status.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Update status likes
		err = s.InsertStatusLikePair(authUser.ID, statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(s, *poster, *status, authUser))
	}
}

// UnikeStatus unlikes status
func UnikeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		statusID := c.Param("statusID")
		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		poster, err := s.GetUserByID(status.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Update status likes
		err = s.DeleteStatusLikePair(authUser.ID, statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildStatusResponse(s, *poster, *status, authUser))
	}
}

// DeleteStatus updates status
func DeleteStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		statusID := c.Param("statusID")
		status, err := s.GetStatusByID(statusID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check if user is poster
		if status.UserID != authUser.ID {
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

		c.JSON(http.StatusOK, buildStatusResponse(s, *authUser, *status, authUser))
	}
}

// GetComments gets comments on a status
func GetComments() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, _ := c.MustGet("authUser").(*store.User)

		statusID := c.Param("statusID")

		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			offset = 0
		} else if offset < 0 {
			offset = 0
		}

		comments, err := s.GetCommentsFromStatus(statusID, 30, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO: do join instead??? or make client find user info instead??
		apiComments := make([]interface{}, len(comments))
		for i, comment := range comments {

			poster, err := s.GetUserByID(comment.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			apiComments[i] = buildCommentResponse(s, *poster, comment, authUser)
		}

		c.JSON(http.StatusOK, gin.H{
			"comments": apiComments,
			"offset":   offset + 30,
		})
	}
}

// PostComment posts a comment on a status
func PostComment() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		commentID := uuid.New().String()
		content := c.PostForm("content")
		statusID := c.Param("statusID")

		comment := store.Comment{
			ID:              commentID,
			UserID:          authUser.ID,
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

		c.JSON(http.StatusOK, buildCommentResponse(s, *authUser, comment, authUser))
	}
}
