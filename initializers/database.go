package initializers

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectToDB() {

	dsn := "host=bubble.db.elephantsql.com user=yktnqwwu password=SGu3XLxLVnJu0-LKhr2nDourm7t8UnfS dbname=yktnqwwu port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed connect to database")
	}
}
