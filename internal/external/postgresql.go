package external

import (
	"fmt"
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
		&models.RecodePowermeter{},
		&models.UserDevice{},
	)
	if err != nil {
		return err
	}

	err = CreateAdminUser(db, cfg)
	if err != nil {
		log.Errorf("failed to create admin user: %v", err)
	}

	err = CreateDeviceMaster(db, cfg)
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
		Name:     cfg.AdminUser.Name,
		Password: string(hashedPassword),
		Role:     models.ADMIN,
	}
	return db.Create(&admin).Error
}

func CreateDeviceMaster(db *gorm.DB, cfg *config.Config) error {
	device := []models.DeviceMaster{
		{ID: cfg.Devices.DEVICE01.DeviceId, Name: cfg.Devices.DEVICE01.Name},
		{ID: cfg.Devices.DEVICE02.DeviceId, Name: cfg.Devices.DEVICE02.Name},
		{ID: cfg.Devices.DEVICE03.DeviceId, Name: cfg.Devices.DEVICE03.Name},
	}
	return db.Create(&device).Error
}
