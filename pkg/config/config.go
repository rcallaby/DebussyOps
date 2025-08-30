package config

import (
	"os"
)

// Minimal config loader - replace with viper for production

type Config struct {
	CalendarURL string
	TodoURL     string
}

func Load() *Config {
	c := &Config{
		CalendarURL: getEnv("CALENDAR_URL", "http://localhost:8081"),
		TodoURL:     getEnv("TODO_URL", "http://localhost:8082"),
	}
	return c
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}