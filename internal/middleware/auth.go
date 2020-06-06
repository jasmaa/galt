package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthUser authorizes user
func AuthUser(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Extract token
		authHeader := c.GetHeader("Authorization")
		r := regexp.MustCompile(`Bearer (.+)`)
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

		// TODO: reply with other errors
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
		}
	}
}
