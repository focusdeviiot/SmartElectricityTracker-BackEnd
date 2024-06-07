package services

import (
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
)

type MasterDeviceService struct {
	deviceRepo *repositories.DeviceRepository
}

func NewMasterDeviceService(deviceRepo *repositories.DeviceRepository) *MasterDeviceService {
	return &MasterDeviceService{
		deviceRepo: deviceRepo,
	}
}

func (s *MasterDeviceService) GetAllDevices() ([]models.DeviceMaster, error) {
	// Get all devices from database
	devices, err := s.deviceRepo.GetAllDevices()
	if err != nil {
		return nil, err
	}

	return devices, nil
}
