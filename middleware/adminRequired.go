package middleware

import (
	"net/http"
	"rental-mobil/initializers"
	"rental-mobil/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You have no access"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserOwnsOrderOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")

		var order models.Order
		if err := initializers.DB.First(&order, "order_id = ?", orderID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		claims, _ := c.MustGet("claims").(*jwt.MapClaims)

		role, _ := c.Get("role")

		if role != "admin" && (order.AdminID != nil && order.AdminID == nil) {
			userIDStr, ok := (*claims)["user_id"].(string)
			if !ok || order.UserID.String() != userIDStr {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this order"})
				c.Abort()
				return
			}
		} else {
			if order.UserID != nil && order.AdminID == nil {
				adminIDStr, ok := (*claims)["admin_id"].(string)
				if ok {
					adminID, err := uuid.Parse(adminIDStr)
					if err == nil {
						order.AdminID = &adminID
						initializers.DB.Save(&order)
					}
				}
			}
		}

		c.Next()
	}
}
