package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL                string
	AppHost              string
	AppPort              string
	JWTSecret            string
	JWTExpirationMinutes int
	BaseURL              string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	appHost := os.Getenv("APP_HOST")

	appPort := os.Getenv("APP_PORT")

	jwtSecret := os.Getenv("JWT_SECRET")

	jwtExpiration := 30
	if raw := os.Getenv("JWT_EXPIRATION_MINUTES"); raw != "" {
		if minutes, err := strconv.Atoi(raw); err == nil && minutes > 0 {
			jwtExpiration = minutes
		}
	}

	dbURL := os.Getenv("DB_URL")

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &Config{
		DBURL:                dbURL,
		AppHost:              appHost,
		AppPort:              appPort,
		JWTSecret:            jwtSecret,
		JWTExpirationMinutes: jwtExpiration,
		BaseURL:              baseURL,
	}
}
