package database

import (
	"fmt"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	// err := CreateUuidOssp(db)
	// if err != nil {
	// 	return err
	// }
	err := CreateUserRoleEnumIfNotExists(db)
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)
}

func CreateUserRoleEnumIfNotExists(db *gorm.DB) error {
	// เช็คว่า enum user_role มีอยู่แล้วหรือไม่
	var exists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role')").Scan(&exists)
	if !exists {
		err := db.Exec("CREATE TYPE user_role AS ENUM ('USER', 'ADMIN')").Error
		if err != nil {
			return fmt.Errorf("failed to create enum user_role: %v", err)
		}
	}
	return nil
}

func CreateUuidOssp(db *gorm.DB) error {
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		return fmt.Errorf("failed to create extension uuid-ossp: %v", err)
	}
	return nil
}
