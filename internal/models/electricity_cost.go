package models

import (
	"time"

	"gorm.io/gorm"
)

type ElectricityCost struct {
	gorm.Model
	ID            string         `json:"id" gorm:"type:varchar(255);default:gen_random_uuid();primary_key;index:;"`
	DeviceID      string         `json:"device_id" gorm:"type:varchar(255);index:;not null;"`
	Device        DeviceMaster   `json:"device" gorm:"foreignKey:DeviceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Cost          int            `json:"total_cost" gorm:"type:integer;index:;not null;"`
	CreatedBy     string         `json:"created_by" gorm:"type:varchar(255);index:;not null;"`
	CreatedByName string         `json:"created_by_name" gorm:"type:varchar(255);index:;not null;"`
	UpdateBy      string         `json:"update_by" gorm:"type:varchar(255);index:;not null;"`
	UpdateByName  string         `json:"update_by_name" gorm:"type:varchar(255);index:;not null;"`
	CreatedAt     time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
