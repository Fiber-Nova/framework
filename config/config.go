package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

// AppConfig holds app-specific configuration
type AppConfig struct {
	Name  string
	Env   string
	Debug bool
	Key   string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver string
	Host   string
	Port   string
	Name   string
	User   string
	Pass   string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "FiberNova"),
			Env:   getEnv("APP_ENV", "development"),
			Debug: getEnvBool("APP_DEBUG", true),
			Key:   getEnv("APP_KEY", ""),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "postgres"),
			Host:   getEnv("DB_HOST", "localhost"),
			Port:   getEnv("DB_PORT", "5432"),
			Name:   getEnv("DB_NAME", "fibernova"),
			User:   getEnv("DB_USER", "root"),
			Pass:   getEnv("DB_PASS", ""),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "3000"),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool retrieves a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}
