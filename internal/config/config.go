package config

import "os"

type Config struct {
	Port           string
	Service        string
	Version        string
	InternalAPIKey string
	DatabaseURL    string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	internalAPIKey := os.Getenv("INTERNAL_API_KEY")
	if internalAPIKey == "" {
		internalAPIKey = "dev-token"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://motor_fiscal:motor_fiscal@localhost:5433/motor_fiscal?sslmode=disable"
	}

	return Config{
		Port:           port,
		Service:        "motor-fiscal",
		Version:        "1.0.0",
		InternalAPIKey: internalAPIKey,
		DatabaseURL:    databaseURL,
	}
}
