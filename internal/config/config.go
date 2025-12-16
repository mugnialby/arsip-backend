package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all environment-based configuration
type Config struct {
	AppName string
	AppEnv  string
	Port    string

	// Database config (used by appcontext)
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Gorm
	AutoMigrate string

	// JWT config
	JWTSecret    string
	JWTExpiresIn int
}

// Load reads configuration from .env or environment variables
func Load() *Config {
	// Load .env file (only for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Parse JWT expiration
	jwtExpStr := getEnv("JWT_EXPIRATION_MINUTES", "60")
	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		jwtExp = 60
	}

	return &Config{
		AppName: getEnv("APP_NAME", "Perpustakaan Backend"),
		AppEnv:  getEnv("APP_ENV", "dev"),
		Port:    getEnv("PORT", "8080"),

		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "admin"),
		DBName:      getEnv("DB_NAME", "db_arsip"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		AutoMigrate: getEnv("DB_AUTO_MIGRATE", "false"),

		JWTSecret:    getEnv("JWT_SECRET", "changeme"),
		JWTExpiresIn: jwtExp,
	}
}

// Helper to get env var or fallback
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
