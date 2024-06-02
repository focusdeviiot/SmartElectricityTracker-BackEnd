package database

import (
	"fmt"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB, cfg *config.Config) error {
	err := CreateUserRoleEnumIfNotExists(db)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.DeviceMaster{},
		&models.UserDevice{},
	)
	if err != nil {
		return err
	}

	err = CreateAdminUser(db, cfg)
	if err != nil {
		log.Errorf("failed to create admin user: %v", err)
	}

	err = CreateDeviceMaster(db)
	if err != nil {
		log.Errorf("failed to create device master: %v", err)
	}

	return nil
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

func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.User{},
		&models.RefreshToken{},
		&models.DeviceMaster{},
		&models.UserDevice{},
	)
}

func CreateAdminUser(db *gorm.DB, cfg *config.Config) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Username: cfg.AdminUser.Username,
		Password: string(hashedPassword),
		Role:     models.ADMIN,
	}
	return db.Create(&admin).Error
}

func CreateDeviceMaster(db *gorm.DB) error {
	device := []models.DeviceMaster{{ID: "DEVICE-01", Name: "DEVICE-01"}, {ID: "DEVICE-02", Name: "DEVICE-02"}, {ID: "DEVICE-03", Name: "DEVICE-03"}}
	return db.Create(&device).Error
}
