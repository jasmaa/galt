package store

// User is a site user
type User struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Username string `form:"username" json:"username"binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
