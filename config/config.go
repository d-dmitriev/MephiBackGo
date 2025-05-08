package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	JWTTtl     string
}

func LoadConfig() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "bankuser"),
		DBPassword: getEnv("DB_PASSWORD", "bankpass"),
		DBName:     getEnv("DB_NAME", "bankdb"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
		JWTTtl:     getEnv("JWT_TTL", "24h"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
