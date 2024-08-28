package controllers

import (
	"log"
	"net/http"
	"os"
	"rental-mobil/initializers"
	"rental-mobil/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAdmin(c *gin.Context) {
	// get data
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if admin admin already exist
	var admin models.Admin
	if err := initializers.DB.Where("email = ?", body.Email).First(&admin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create the admin
	admin = models.Admin{
		AdminID:  uuid.New(),
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	if err := initializers.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Return the admin
	c.JSON(http.StatusOK, gin.H{"admin": admin})
}

func GenerateTokenAdmin(admin models.Admin) (string, error) {
	var err error
	secret := os.Getenv("JWT_SECRET_KEY")

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.AdminID,
		"email":    admin.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"role":     "admin",
	})

	// Sign the token with secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func LoginAdmin(c *gin.Context) {
	var admin models.Admin
	var inputAdmin models.Admin

	// Bind the admin input to inputAdmin
	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find admin by email
	if err := initializers.DB.Where("email = ?", inputAdmin.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please check your username/email details"})
		return
	}

	// Compare the password with the hashed password stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(inputAdmin.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please check your password details"})
		return
	}

	// If the password is correct, generate a token for the admin
	token, err := GenerateTokenAdmin(admin)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
