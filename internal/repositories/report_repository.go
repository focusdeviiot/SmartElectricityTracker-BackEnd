package repositories

import (
	"smart_electricity_tracker_backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) FindReportByDeviceAndDate(device_id *string, dateFrom *time.Time, dateTo *time.Time) ([]models.ReportRes, error) {
	var reports []models.ReportRes
	query := r.db.Table("recode_powermeters as rp")
	if device_id != nil && *device_id != "" {
		query = query.Where("rp.device_id = ?", device_id)
	}
	if dateFrom != nil {
		query = query.Where("rp.created_at >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("rp.created_at <= ?", dateTo)
	}
	if err := query.Select("rp.id, rp.device_id, rp.volt, rp.ampere, rp.watt, rp.created_at").Scan(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportRepository) RecordPowermeter(report *models.RecodePowermeter) error {
	if err := r.db.Create(report).Error; err != nil {
		return err
	}
	return nil
}
