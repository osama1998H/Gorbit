// internal/database/postgres.go
package database

import (
	"fmt"
	"log"

	"gorbit/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Databases.Postgres.Host,
		cfg.Databases.Postgres.Username,
		cfg.Databases.Postgres.Password,
		cfg.Databases.Postgres.Database,
		cfg.Databases.Postgres.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Add additional configurations if needed
	})
	if err != nil {
		log.Printf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	return db, nil
}
