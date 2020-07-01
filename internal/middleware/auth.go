package middleware

import (
	"errors"
	"net/http"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/jasmaa/galt/internal/store"
)

// AuthUser authorizes user
func AuthUser(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Extract token
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) == 0 {
			c.Set("authUser", nil)
			c.Next()
			return
		}

		r := regexp.MustCompile(`Bearer ([\w-]+\.[\w-]+\.[\w-]+)`)
		res := r.FindStringSubmatch(authHeader)
		if len(res) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		tokenString := res[1]

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			// Validate alg
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}

			return []byte(hmacSecret), nil
		})

		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			// Get authenticated user
			s := c.MustGet("store").(store.Store)
			authUserID := claims["userID"].(string)
			authUser, err := s.GetUserByID(authUserID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token",
				})
				c.Abort()
				return
			}

			c.Set("authUser", authUser)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
		}
	}
}
