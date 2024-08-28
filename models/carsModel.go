package models

import (
	"gorm.io/gorm"
)

type Car struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100)"`
	CarType   string `gorm:"type:varchar(50)"`
	Rating    float32
	Fuel      string `gorm:"type:varchar(50)"`
	Image     string `gorm:"type:text"`
	HourRate  float32
	DayRate   float32
	MonthRate float32
}
