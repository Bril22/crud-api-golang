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

func RegisterUsers(c *gin.Context) {
	// get data
	var body struct {
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Username    string `json:"username"`
		City        string `json:"city"`
		Zip         string `json:"zip"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists
	var user models.User
	if err := initializers.DB.Where("email = ?", body.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create the user
	user = models.User{
		UserID:      uuid.New(),
		Email:       body.Email,
		Password:    string(hashedPassword),
		Username:    body.Username,
		City:        body.City,
		Zip:         body.Zip,
		PhoneNumber: body.PhoneNumber,
		Address:     string(body.Address),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Return the user
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GenerateTokenUser(user models.User) (string, error) {
	var err error
	secret := os.Getenv("JWT_SECRET_KEY")

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"role":    "user",
	})

	// Sign the token with secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func LoginUser(c *gin.Context) {
	var user models.User
	var inputUser models.User

	// Bind the user input to inputUser
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	if err := initializers.DB.Where("email = ?", inputUser.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please check your username/email details"})
		return
	}

	// Compare the password with the hashed password stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please check your password details"})
		return
	}

	// If the password is correct, generate a token for the user
	token, err := GenerateTokenUser(user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
