package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	var secretKey = os.Getenv("JWT_SECRET_KEY")
	return func(c *gin.Context) {
		// Get the client secret key
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Verify the token
		claims := &jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(extractedToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		// cek log
		fmt.Println("Parsed Token:", parsedToken)
		fmt.Println("Claims:", claims)

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing token"})
			}
			c.Abort()
			return
		}

		if parsedToken.Valid {
			if _, ok := (*claims)["admin_id"]; ok {
				c.Set("role", "admin")
			} else {
				c.Set("role", "user")
			}

			c.Set("claims", claims)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
