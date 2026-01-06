package appcontext

import (
	"fmt"

	"github.com/mugnialby/arsip-backend/internal/config"
	"github.com/mugnialby/arsip-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppContext struct {
	DB *gorm.DB
}

func NewAppContext(cfg *config.Config) (*AppContext, error) {
	// Build PostgreSQL DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// --- AUTO MIGRATION ---
	// This will create or update the schema automatically for both tables.
	// autoMigrate, _ := strconv.ParseBool(cfg.AutoMigrate)
	// if autoMigrate {
	// 	db.AutoMigrate(
	// 		&model.Category{},
	// 		&model.Archive{},
	// 	)
	// }

	logger.Log.Info("main.context.database.success",
		zap.Any("message", "Database connected and migrated successfully"),
	)

	return &AppContext{
		DB: db,
	}, nil
}
