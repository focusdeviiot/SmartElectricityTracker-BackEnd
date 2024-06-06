package services

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"sync"
	"time"
)

type ReportService struct {
	reportRepo *repositories.ReportRepository
	mu         *sync.Mutex
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

func (s *ReportService) RecordPowermeter(report *models.RecodePowermeter) error {
	return s.reportRepo.RecordPowermeter(report)
}
