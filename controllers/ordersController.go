package controllers

import (
	"errors"
	"net/http"
	"rental-mobil/initializers"
	"rental-mobil/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func PostOrder(c *gin.Context) {
	var body struct {
		PickUpLoc   string    `json:"pick_up_loc" binding:"required"`
		DropOffLoc  string    `json:"drop_off_loc" binding:"required"`
		PickUpDate  time.Time `json:"pick_up_date"`
		DropOffDate time.Time `json:"drop_off_date"`
		PickUpTime  time.Time `json:"pick_up_time"`
		CarID       uint      `json:"car_id" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var carCount int64
	initializers.DB.Model(&models.Car{}).Where("id = ?", body.CarID).Count(&carCount)
	if carCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CarID does not exist"})
		return
	}

	claims := c.MustGet("claims").(*jwt.MapClaims)

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		c.Abort()
		return
	}

	order := models.Order{
		PickUpLoc:   body.PickUpLoc,
		DropOffLoc:  body.DropOffLoc,
		PickUpDate:  body.PickUpDate,
		DropOffDate: body.DropOffDate,
		PickUpTime:  body.PickUpTime,
		CarID:       body.CarID,
	}

	if role == "admin" {
		adminIDStr, ok := (*claims)["admin_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin_id in claims"})
			return
		}
		adminID, err := uuid.Parse(adminIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin_id format"})
			return
		}
		order.AdminID = &adminID
		order.UserID = nil
	} else {
		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in claims"})
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id format"})
			return
		}
		order.UserID = &userID
		order.AdminID = nil
	}

	result := initializers.DB.Create(&order)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	// return post
	c.JSON(200, gin.H{
		"orders": order,
	})
}

func GetListOrders(c *gin.Context) {
	var orders []models.Order

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		c.Abort()
		return
	}

	claims := c.MustGet("claims").(*jwt.MapClaims)
	if role == "admin" {
		initializers.DB.Order("id asc").Find(&orders)
	} else {
		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in claims"})
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id format"})
			return
		}
		initializers.DB.Where("user_id = ?", userID).Order("id asc").Find(&orders)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func UpdateDataOrder(c *gin.Context) {
	orderID := c.Param("id")

	var body struct {
		PickUpLoc   string    `json:"pick_up_loc" binding:"required"`
		DropOffLoc  string    `json:"drop_off_loc" binding:"required"`
		PickUpDate  time.Time `json:"pick_up_date"`
		DropOffDate time.Time `json:"drop_off_date"`
		PickUpTime  time.Time `json:"pick_up_time"`
		CarID       uint      `json:"car_id" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var carCount int64
	initializers.DB.Model(&models.Car{}).Where("id = ?", body.CarID).Count(&carCount)
	if carCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CarID does not exist"})
		return
	}

	claims := c.MustGet("claims").(*jwt.MapClaims)

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		c.Abort()
		return
	}

	order := models.Order{
		PickUpLoc:   body.PickUpLoc,
		DropOffLoc:  body.DropOffLoc,
		PickUpDate:  body.PickUpDate,
		DropOffDate: body.DropOffDate,
		PickUpTime:  body.PickUpTime,
		CarID:       body.CarID,
	}

	if role == "admin" {
		adminIDStr, ok := (*claims)["admin_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin_id in claims"})
			return
		}
		adminID, err := uuid.Parse(adminIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin_id format"})
			return
		}
		order.AdminID = &adminID
		order.UserID = nil
	} else {
		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in claims"})
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id format"})
			return
		}
		order.UserID = &userID
		order.AdminID = nil
	}

	result := initializers.DB.First(&order, "order_id = ?", orderID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		// Handle other types of errors if necessary
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	results := initializers.DB.Model(&order).Updates(models.Order{
		PickUpLoc:   body.PickUpLoc,
		DropOffLoc:  body.DropOffLoc,
		PickUpDate:  body.PickUpDate,
		DropOffDate: body.DropOffDate,
		PickUpTime:  body.PickUpTime,
		CarID:       body.CarID,
	})
	if results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": order,
	})
}

func DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Fetch the existing order first
	var order models.Order
	result := initializers.DB.First(&order, "order_id = ?", orderID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		c.Abort()
		return
	}

	claims := c.MustGet("claims").(*jwt.MapClaims)
	if role != "admin" {
		userIDStr, ok := (*claims)["user_id"].(string)
		if order.UserID == nil || !ok || order.UserID.String() != userIDStr {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this order"})
			return
		}
	}

	// Delete the order
	initializers.DB.Delete(&order)

	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted successfully",
	})
}
