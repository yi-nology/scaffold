package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"scaffold/internal/conf"
)

var DB *gorm.DB

func Init(cfg *conf.DatabaseConfig) error {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN())
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	case "sqlite":
		// Auto-create parent directory for SQLite file
		if dir := filepath.Dir(cfg.DBName); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create sqlite directory: %w", err)
			}
		}
		dialector = sqlite.Open(cfg.DBName)
	default:
		return fmt.Errorf("unsupported database driver: %s (supported: mysql, postgres, sqlite)", cfg.Driver)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		DB = nil // ensure DB is nil on failure
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return nil
	}
	return DB.AutoMigrate(models...)
}
