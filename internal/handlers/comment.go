package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jasmaa/galt/internal/store"
)

// APIComment represents a status in API response
type APIComment struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          APIUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

// TODO: figure out response for auth vs not auth
// buildCommentResponse builds API comment response
func buildCommentResponse(poster store.User, comment store.Comment, commentLikes int) APIComment {
	return APIComment{
		ID: comment.ID,
		Poster: APIUser{
			ID:            poster.ID,
			Username:      poster.Username,
			Description:   poster.Description,
			ProfileImgURL: poster.ProfileImgURL,
		},
		Content:         comment.Content,
		Likes:           commentLikes,
		PostedTimestamp: comment.PostedTimestamp,
		IsEdited:        comment.IsEdited,
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
		apiComments := make(map[string]APIComment)
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
		userID := c.MustGet("userID").(string)
		content := c.PostForm("content")
		statusID := c.Param("statusID")

		comment := store.Comment{
			ID:              commentID,
			UserID:          userID,
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

		user, err := s.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCommentResponse(*user, comment, 0))
	}
}
