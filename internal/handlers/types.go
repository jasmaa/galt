package handlers

import (
	"time"

	"github.com/jasmaa/galt/internal/store"
)

type apiUser struct {
	ID            string `form:"id" json:"id" binding:"required"`
	Username      string `form:"username" json:"username" binding:"required"`
	Description   string `form:"description" json:"description" binding:"required"`
	ProfileImgURL string `form:"profileImgURL" json:"profileImgURL" binding:"required"`
}

func buildUserResponse(user store.User) apiUser {
	return apiUser{
		ID:            user.ID,
		Username:      user.Username,
		Description:   user.Description,
		ProfileImgURL: user.ProfileImgURL,
	}
}

type apiStatusNonAuth struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          apiUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	Reshares        int       `form:"reshares" json:"reshares" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

type apiStatusAuth struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          apiUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	IsLiked         bool      `form:"isLiked" json:"isLiked" binding:"required"`
	Reshares        int       `form:"reshares" json:"reshares" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

func buildStatusResponse(s store.Store, poster store.User, status store.Status, authUser *store.User) interface{} {

	statusLikes, _ := s.GetStatusLikes(status.ID)

	// Return authenticated response
	if authUser != nil {

		isLiked, _ := s.GetIsUserLikedStatus(authUser.ID, status.ID)

		return apiStatusAuth{
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
			Reshares:        -2, // filler value for now
			PostedTimestamp: status.PostedTimestamp,
			IsEdited:        status.IsEdited,
		}
	}

	return apiStatusNonAuth{
		ID: status.ID,
		Poster: apiUser{
			ID:            poster.ID,
			Username:      poster.Username,
			Description:   poster.Description,
			ProfileImgURL: poster.ProfileImgURL,
		},
		Content:         status.Content,
		Likes:           statusLikes,
		Reshares:        -1, // filler value for now
		PostedTimestamp: status.PostedTimestamp,
		IsEdited:        status.IsEdited,
	}
}

type apiCommentNonAuth struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          apiUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

type apiCommentAuth struct {
	ID              string    `form:"id" json:"id" binding:"required"`
	Poster          apiUser   `form:"poster" json:"poster" binding:"required"`
	Content         string    `form:"content" json:"content" binding:"required"`
	Likes           int       `form:"likes" json:"likes" binding:"required"`
	IsLiked         bool      `form:"isLiked" json:"isLiked" binding:"required"`
	PostedTimestamp time.Time `form:"postedTimestamp" json:"postedTimestamp" binding:"required"`
	IsEdited        bool      `form:"isEdited" json:"isEdited" binding:"required"`
}

func buildCommentResponse(s store.Store, poster store.User, comment store.Comment, authUser *store.User) interface{} {

	commentLikes, _ := s.GetCommentLikes(comment.ID)

	if authUser != nil {
		return apiCommentAuth{
			ID: comment.ID,
			Poster: apiUser{
				ID:            poster.ID,
				Username:      poster.Username,
				Description:   poster.Description,
				ProfileImgURL: poster.ProfileImgURL,
			},
			Content:         comment.Content,
			Likes:           commentLikes,
			IsLiked:         false, // filler value
			PostedTimestamp: comment.PostedTimestamp,
			IsEdited:        comment.IsEdited,
		}
	}

	return apiCommentNonAuth{
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

type apiCircle struct {
	ID          string `form:"id" json:"id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	UserCount   int    `form:"userCount" json:"userCount" binding:"required"`
}

func buildCircleResponse(s store.Store, circle store.Circle) apiCircle {

	userCount, _ := s.GetCircleUserCount(circle.ID)

	return apiCircle{
		ID:          circle.ID,
		Name:        circle.Name,
		Description: circle.Description,
		UserCount:   userCount,
	}
}
