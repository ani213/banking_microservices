package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthService   string
	SMTPHost      string
	SMTPPort      string
	EmailPassword string
	EmailId       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
	return &Config{
		AuthService:   os.Getenv("AUTH_SERVICE"),
		SMTPHost:      os.Getenv("SMTPHOST"),
		SMTPPort:      os.Getenv("SMTPPORT"),
		EmailPassword: os.Getenv("PASSWORD"),
		EmailId:       os.Getenv("EMAIL"),
	}
}
