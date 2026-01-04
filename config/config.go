package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port        string
	Environment string
	LogLevel    string

	SubscanBaseURL string
	SubscanAPIKey  string

	MaxPages    int
	RowsPerPage int

	HTTPTimeoutSeconds int

	BlocksPerEra int64
	PlanckToDOT  float64
}

var cfg *Config

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		logrus.Debug("No .env file found, using environment variables")
	}

	cfg = &Config{}

	cfg.Port = getEnvOrDefault("PORT", "8080")
	cfg.Environment = getEnvOrDefault("ENVIRONMENT", "development")
	cfg.LogLevel = getEnvOrDefault("LOG_LEVEL", "info")
    
	cfg.SubscanBaseURL = getEnvOrDefault("SUBSCAN_BASE_URL", "https://polkadot.api.subscan.io")
	cfg.SubscanAPIKey = getEnvOrDefault("SUBSCAN_API_KEY", "")

	cfg.MaxPages = getEnvAsIntOrDefault("MAX_PAGES", 10)
	cfg.RowsPerPage = getEnvAsIntOrDefault("ROWS_PER_PAGE", 100)

	cfg.HTTPTimeoutSeconds = getEnvAsIntOrDefault("HTTP_TIMEOUT_SECONDS", 30)

	cfg.BlocksPerEra = int64(getEnvAsIntOrDefault("BLOCKS_PER_ERA", 14400))
	cfg.PlanckToDOT = getEnvAsFloatOrDefault("PLANCK_TO_DOT", 1e10)

	return cfg
}

func Get() *Config {
	if cfg == nil {
		logrus.Warnf("Configuration not loaded. Call config.Load() first")
	}
	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		logrus.Warnf("Invalid integer value for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
}

func getEnvAsFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
		logrus.Warnf("Invalid float value for %s, using default: %f", key, defaultValue)
	}
	return defaultValue
}
