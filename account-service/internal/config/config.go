package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource  string
	JwtSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}
	return &Config{
		DBSource:  os.Getenv("DATABASE_URL"),
		JwtSecret: os.Getenv("SECRET_KEY"),
	}

}
