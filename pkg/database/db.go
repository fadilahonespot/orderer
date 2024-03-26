package database

import (
	"fmt"
	"os"

	"github.com/fadilahonespot/mygram/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	fmt.Println("Connect to Database")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if os.Getenv("DB_DEBUG") == "true" {
		db = db.Debug()
	}

	db.AutoMigrate(
		&entity.User{},
		&entity.Photo{},
		&entity.Comment{},
		&entity.SocialMedia{},
	)

	return db
}
