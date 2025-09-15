package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"going/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
	dbErr  error

	// models to be registered
	models = make([]interface{}, 0)

	// ErrNotConnected is returned when the database is not connected
	ErrNotConnected = errors.New("database not connected")
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*sql.DB, error) {
	dbOnce.Do(func() {
		switch cfg.Database.Driver {
		case "sqlite3":
			initSQLite(cfg)
		default:
			dbErr = fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
		}
	})

	if dbErr != nil {
		return nil, dbErr
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return sqlDB, nil
}

// initSQLite initializes a SQLite database
func initSQLite(cfg *config.Config) {
	// Ensure the database directory exists
	if err := os.MkdirAll(cfg.Database.Path, 0755); err != nil {
		dbErr = fmt.Errorf("failed to create database directory: %w", err)
		return
	}

	dbPath := filepath.Join(cfg.Database.Path, cfg.Database.Name)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to the database
	db, dbErr = gorm.Open(sqlite.Open(dbPath), gormConfig)
	if dbErr != nil {
		dbErr = fmt.Errorf("failed to connect to database: %w", dbErr)
	}
}

// GetDB returns the database instance
func GetDB() (*gorm.DB, error) {
	if db == nil {
		return nil, ErrNotConnected
	}
	return db, nil
}

// RegisterModels registers models for auto-migration
func RegisterModels(modelList ...interface{}) {
	models = append(models, modelList...)
}

// runMigrations runs database migrations
func runMigrations() error {
	if db == nil {
		return ErrNotConnected
	}

	// Auto migrate all registered models
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return fmt.Errorf("failed to migrate models: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return sqlDB.Close()
}
