package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	EmailServer string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Env file not loaded")
	}
	return &Config{
		EmailServer: os.Getenv("EMAIL_SERVER"),
	}

}
