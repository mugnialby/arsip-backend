package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mugnialby/arsip-backend/internal/utils"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"go.uber.org/zap"
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

func Load() *Config {
	envFilePath, err := utils.GetEnvFilePath()
	if err != nil {
		logger.Log.Error("config.env.get_env_file_path.failed",
			zap.Error(err),
			zap.String("message", "Failed to get .env file path"),
		)
	}

	if err := godotenv.Load(envFilePath); err != nil {
		logger.Log.Warn("config.env.load.failed",
			zap.Error(err),
			zap.String("path", envFilePath),
			zap.String("message", "No .env file found, using system environment variables"),
		)
	}

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

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
