// internal/database/mysql.go
package database

import (
	"fmt"
	"log"

	"gorbit/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Databases.MySQL.Username,
		cfg.Databases.MySQL.Password,
		cfg.Databases.MySQL.Host,
		cfg.Databases.MySQL.Port,
		cfg.Databases.MySQL.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Add additional configurations if needed
	})
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
		return nil, err
	}

	return db, nil
}
