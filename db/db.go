package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

func SetupDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")

	args := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, name, user, password, sslmode)

	db, err := gorm.Open("postgres", args)
	return db, err
}
