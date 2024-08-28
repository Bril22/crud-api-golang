package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email       string    `gorm:"type:varchar(100);unique_index"`
	Password    string    `gorm:"type:varchar(100)"`
	PhoneNumber string
	City        string
	Zip         string
	Message     string
	Username    string
	Address     string
}

type Admin struct {
	gorm.Model
	AdminID  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email    string    `gorm:"type:varchar(100);unique_index"`
	Password string    `gorm:"type:varchar(100)"`
}

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
