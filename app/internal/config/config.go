package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}

	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" {
		panic("missing required environment variables: DATABASE_URL and JWT_SECRET must be set")
	}

	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg
}
