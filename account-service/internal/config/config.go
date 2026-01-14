package config

import (
	"log"
	"os"

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
	conn, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
	if err != nil {
		log.Println("Rabbitmq is not conneted", err.Error())
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
