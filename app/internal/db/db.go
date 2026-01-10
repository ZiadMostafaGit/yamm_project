package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10) //becase its not production if i was to make it production i would increasee this
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
