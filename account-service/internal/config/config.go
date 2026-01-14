package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	DBSource     string
	JwtSecret    string
	EmailServer  string
	QueueChannel *amqp.Channel
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
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
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Creating channel getting issue", err.Error())
	}
	return &Config{
		DBSource:     os.Getenv("DATABASE_URL"),
		JwtSecret:    os.Getenv("SECRET_KEY"),
		EmailServer:  os.Getenv("EMAIL_SERVER"),
		QueueChannel: ch,
	}

}
