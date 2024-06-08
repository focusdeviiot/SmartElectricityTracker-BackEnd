package services

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"time"
)

type ReportService struct {
	reportRepo *repositories.ReportRepository
	cfg        *config.Config
}

func NewReportService(reportRepo *repositories.ReportRepository, cfg *config.Config) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		cfg:        cfg,
	}
}

func (s *ReportService) GetReportByDeviceAndDate(device_id *string, dateFrom *string, dateTo *string) ([]models.ReportRes, error) {
	dateFromSet, err := time.Parse(time.RFC3339, *dateFrom)
	if err != nil {
		return nil, err
	}
	dateToSet, err := time.Parse(time.RFC3339, *dateTo)
	if err != nil {
		return nil, err
	}
	return s.reportRepo.FindReportByDeviceAndDate(device_id, &dateFromSet, &dateToSet)
}
