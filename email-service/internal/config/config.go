package config

import (
	"log"
	"os"

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
	conn, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
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
