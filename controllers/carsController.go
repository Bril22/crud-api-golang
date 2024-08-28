package controllers

import (
	"errors"
	"net/http"
	"rental-mobil/initializers"
	"rental-mobil/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostCar(c *gin.Context) {
	var body struct {
		Name      string  `json:"name" binding:"required"`
		CarType   string  `json:"car_type" binding:"required"`
		Rating    float32 `json:"rating" binding:"required"`
		Fuel      string  `json:"fuel" binding:"required"`
		Image     string  `json:"image"`
		HourRate  float32 `json:"hour_rate" binding:"required"`
		DayRate   float32 `json:"day_rate" binding:"required"`
		MonthRate float32 `json:"month_rate" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	car := models.Car{Name: body.Name, CarType: body.CarType, Rating: body.Rating, Fuel: body.Fuel, Image: body.Image, HourRate: body.HourRate, DayRate: body.DayRate, MonthRate: body.MonthRate}

	result := initializers.DB.Create(&car)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	// return post car
	c.JSON(200, gin.H{
		"cars": car,
	})
}

func GetListCar(c *gin.Context) {
	var cars []models.Car
	initializers.DB.Order("id asc").Find(&cars)
	c.JSON(200, gin.H{
		"cars": cars,
	})
}

func GetListCarByID(c *gin.Context) {
	carID := c.Param("id")
	// get post
	var car models.Car
	initializers.DB.First(&car, carID)

	// response
	c.JSON(200, gin.H{
		"cars": car,
	})
}

func UpdateDataCar(c *gin.Context) {
	carID := c.Param("id")

	var car models.Car
	result := initializers.DB.First(&car, "id = ?", carID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
			return
		}
		// Handle other types of errors if necessary
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var body models.Car
	if err := c.Bind(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	initializers.DB.Model(&car).Updates(models.Car{
		Name:      body.Name,
		CarType:   body.CarType,
		Rating:    body.Rating,
		Fuel:      body.Fuel,
		Image:     body.Image,
		HourRate:  body.HourRate,
		DayRate:   body.DayRate,
		MonthRate: body.MonthRate,
	})

	c.JSON(200, gin.H{
		"car": car,
	})
}

func DeleteCar(c *gin.Context) {
	carID := c.Param("id")

	var car models.Car
	if err := initializers.DB.First(&car, carID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Error delete the Car",
		})
		return
	}

	// delete
	initializers.DB.Delete(&models.Car{}, carID)

	// response
	c.JSON(200, gin.H{
		"message": "Success Delete the Car",
	})
	return
}
