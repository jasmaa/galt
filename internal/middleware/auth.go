package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthUser authorizes user
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)

		if session.Get("username") == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
		}

		c.Next()
	}
}
