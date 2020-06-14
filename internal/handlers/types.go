package handlers

import (
	"time"

	"github.com/jasmaa/galt/internal/store"
)

// APIUser represents a user in API response
type apiUser struct {
	ID            string `form:"id" json:"id" binding:"required"`
	Username      string `form:"username" json:"username" binding:"required"`
	Description   string `form:"description" json:"description" binding:"required"`
	ProfileImgURL string `form:"profileImgURL" json:"profileImgURL" binding:"required"`
}

// buildUserResponse builds API user response
func buildUserResponse(user store.User) apiUser {
	return apiUser{
		ID:            user.ID,
		Username:      user.Username,
		Description:   user.Description,
		ProfileImgURL: user.ProfileImgURL,
	}
}

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

// apiComment represents a status in API response
type apiComment struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          apiUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

// TODO: figure out response for auth vs not auth
// buildCommentResponse builds API comment response
func buildCommentResponse(poster store.User, comment store.Comment, commentLikes int) apiComment {
	return apiComment{
		ID: comment.ID,
		Poster: apiUser{
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

// apiCircle represents a status in API response
type apiCircle struct {
	ID          string `form:"id" json:"id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
}

// buildStatusResponse builds API status response
func buildCircleResponse(circle store.Circle) apiCircle {
	return apiCircle{
		ID:          circle.ID,
		Name:        circle.Name,
		Description: circle.Description,
	}
}
