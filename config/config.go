package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     string
	DBPath   string
	LogLevel string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:     getEnv("PORT", "8006"),
		DBPath:   getEnv("DB_PATH", "./booktracker.db"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	port := c.Port
	if port[0] == ':' {
		port = port[1:]
	}
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("PORT must be a number, got: %s", c.Port)
	}

	if c.DBPath == "" {
		return fmt.Errorf("DB_PATH cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be debug, info, or error, got: %s", c.LogLevel)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
