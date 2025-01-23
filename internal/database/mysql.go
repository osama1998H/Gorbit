// internal/database/mysql.go
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorbit/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Databases.MySQL.Username,
		cfg.Databases.MySQL.Password,
		cfg.Databases.MySQL.Host,
		cfg.Databases.MySQL.Port,
		cfg.Databases.MySQL.Database,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	if cfg.Server.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("mysql connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("mysql pool setup failed: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(cfg.Databases.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Databases.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Databases.MySQL.ConnMaxLifetime) * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("mysql ping failed: %w", err)
	}

	slog.Info("MySQL connection established",
		"host", cfg.Databases.MySQL.Host,
		"database", cfg.Databases.MySQL.Database,
	)

	return db, nil
}
