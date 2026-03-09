package main

import (
	"log"
	"time"

	"bzy/deployer/pkg/lms"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// DB
	db, err := gorm.Open(sqlite.Open("lms.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&lms.LMS{}); err != nil {
		log.Fatal(err)
	}

	// Repository
	repo := lms.NewGormLMSRepository(db)

	// RabbitMQ
	conn, ch, err := connectRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer ch.Close()

	// Service
	service := lms.NewLMSService(repo, ch)

	// LMS directories
	dirs := []string{"C:/inetpub/wwwroot/lms1", "C:/inetpub/wwwroot/lms2"}

	// Start scanner
	lms.RunScanner(service, dirs, 30*time.Second)
}

// connectRabbitMQ helper
func connectRabbitMQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}
