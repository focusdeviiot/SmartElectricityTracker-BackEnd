package repositories

import (
	"smart_electricity_tracker_backend/internal/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) GetAllDevices() ([]models.DeviceMaster, error) {
	// Get all devices from database
	var devices []models.DeviceMaster
	err := r.db.Find(&devices).Where("deleted_at IS NULL")
	if err.Error != nil {
		return nil, err.Error
	}

	return devices, nil
}
