package config

import "os"

type Config struct {
	Port    string
	Service string
	Version string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port:    port,
		Service: "motor-fiscal",
		Version: "1.0.0",
	}
}
