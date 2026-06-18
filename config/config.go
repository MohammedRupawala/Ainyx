package config

import (
	"os"
)

type Config struct {
	DBURL string
	Port  string
}

func Load() *Config {
	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return &Config{
		DBURL: dbURL,
		Port:  port,
	}
}
