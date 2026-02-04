package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	DBHost  string
	DBUser  string
	DBPass  string
	DBName  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

    cfg := &Config{
        AppPort: os.Getenv("APP_PORT"),
        DBHost: os.Getenv("DB_HOST"),
        DBUser: os.Getenv("DB_USER"),
        DBPass: os.Getenv("DB_PASS"),
        DBName: os.Getenv("DB_NAME"),
    }

    return cfg
}