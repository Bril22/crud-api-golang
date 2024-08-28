package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	PickUpLoc   string    `gorm:"type:varchar(100)"`
	DropOffLoc  string    `gorm:"type:varchar(100)"`
	PickUpDate  time.Time
	DropOffDate time.Time
	PickUpTime  time.Time
	CarID       uint
	UserID      *uuid.UUID `gorm:"type:uuid"`
	AdminID     *uuid.UUID `gorm:"type:uuid"`
}
