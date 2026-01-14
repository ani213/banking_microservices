package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	AuthService   string
	SMTPHost      string
	SMTPPort      string
	EmailPassword string
	EmailId       string
	QueueChannel  *amqp.Channel
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
	// conn, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
	amp_url := os.Getenv("RABBIT_MQ_URL")
	var conn *amqp.Connection
	// var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(amp_url)
		if err == nil {
			log.Println("RabbitMQ connected")

		}

		log.Printf("RabbitMQ not ready (attempt %d): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Println(err.Error() + "Error in rabbitmq connetion")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err.Error() + "Error in rabbitmq channel")
	}
	return &Config{
		AuthService:   os.Getenv("AUTH_SERVICE"),
		SMTPHost:      os.Getenv("SMTPHOST"),
		SMTPPort:      os.Getenv("SMTPPORT"),
		EmailPassword: os.Getenv("PASSWORD"),
		EmailId:       os.Getenv("EMAIL"),
		QueueChannel:  ch,
	}
}
