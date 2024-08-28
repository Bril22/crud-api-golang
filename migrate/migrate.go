package main

import (
	"log"
	"rental-mobil/initializers"
	"rental-mobil/models"
)

func init() {
	initializers.ConnectToDB()
	initializers.LoadEnvVariables()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Admin{}, &models.Car{}, &models.Order{})
	if err := initializers.DB.AutoMigrate(&models.User{}, &models.Admin{}, &models.Car{}, &models.Order{}); err != nil {
		log.Fatalf("Error auto-migrating: %v", err)
	}
}
