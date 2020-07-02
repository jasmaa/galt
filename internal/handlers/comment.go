package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jasmaa/galt/internal/store"
)

// GetCommentChain gets a comment chain by parent comment
func GetCommentChain() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		commentID := c.Param("commentID")
		authUser, _ := c.MustGet("authUser").(*store.User)

		c.JSON(http.StatusOK, buildCommentChainResponse(s, commentID, authUser, 3))
	}
}

// PostReply posts a comment reply to another comment
func PostReply() gin.HandlerFunc {
	return func(c *gin.Context) {

		s := c.MustGet("store").(store.Store)
		commentID := c.Param("commentID")
		authUser, ok := c.MustGet("authUser").(*store.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
			return
		}

		replyID := uuid.New().String()
		content := c.PostForm("content")

		parentComment, err := s.GetCommentByID(commentID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		comment := store.Comment{
			ID:              replyID,
			UserID:          authUser.ID,
			StatusID:        parentComment.StatusID,
			ParentCommentID: sql.NullString{String: parentComment.ID, Valid: true},
			Content:         content,
			PostedTimestamp: time.Now(),
			IsEdited:        false,
		}
		err = s.InsertComment(comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, buildCommentResponse(s, *authUser, comment, authUser))
	}
}
