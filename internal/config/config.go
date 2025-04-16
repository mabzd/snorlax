package config

import (
	"os"
)

type Config struct {
	ApiPort            string
	DbHost             string
	DbPort             string
	DbUser             string
	DbPass             string
	DbName             string
	ServerTimeoutInSec int
}

func LoadConfig() Config {
	return Config{
		ApiPort:            getenv("API_PORT", "8080"),
		DbHost:             getenv("DB_HOST", "localhost"),
		DbPort:             getenv("DB_PORT", "5432"),
		DbUser:             getenv("DB_USER", "postgres"),
		DbPass:             getenv("DB_PASS", "postgres"),
		DbName:             getenv("DB_NAME", "snorlax_db"),
		ServerTimeoutInSec: 30,
	}
}

func getenv(key string, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}
	return env
}
