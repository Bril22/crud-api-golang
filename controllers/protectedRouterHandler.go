package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func ProtectedRouteHandler(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if exists {
// 		c.JSON(http.StatusOK, gin.H{
// 			"user":    user.(*jwt.Token).Claims,
// 			"message": "Welcome to the protected route!",
// 		})
// 	}
// }

// func ProtectedRouteHandler(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if exists {
// 		fmt.Println("User Claims:", claims)
// 		c.JSON(http.StatusOK, gin.H{
// 			"user":    claims,
// 			"message": "Welcome to the protected route!",
// 		})
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
// 	}
// }

func ProtectedRouteHandler(c *gin.Context) {
	claims, exists := c.Get("claims")
	if exists {
		role, _ := c.Get("role")
		fmt.Println("User Claims:", claims)
		c.JSON(http.StatusOK, gin.H{
			"user":    claims,
			"role":    role,
			"message": "Welcome to the Rent Car App!",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims or role not found"})
	}
}
