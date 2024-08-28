package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID    int64
	Title string
	Body  string
}
