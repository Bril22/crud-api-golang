package handlers

import (
	"net/http"
	"rental-mobil/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	// Get the order ID from the URL.
	orderID := c.Param("order_id")
	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	var order models.Order

	// Fetch the order and its related User, Admin, and Car information from the database.
	if err := h.db.Preload("User").Preload("Admin").Preload("Car").Find(&order, "order_id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching order"})
		return
	}

	// Return the fetched order.
	c.JSON(http.StatusOK, gin.H{"order": order})
}
