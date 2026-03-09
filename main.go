package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn, err := dsn_from_env()
	_, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
}

func dsn_from_env() (string, error) {

	host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_user_name := os.Getenv("DB_USER_NAME")
	db_user_password := os.Getenv("DB_USER_PASSWORD")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai", host, db_user_name, db_user_password, db_name), nil
}
