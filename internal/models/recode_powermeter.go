package models

import (
	"time"

	"gorm.io/gorm"
)

type RecodePowermeter struct {
	gorm.Model
	ID        int64          `json:"id" gorm:"bigserial; primary_key;"`
	DeviceID  string         `json:"device_id" gorm:"type:varchar(255);index:;not null;"`
	Device    DeviceMaster   `json:"device" gorm:"foreignKey:DeviceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Volt      float32        `json:"volt" gorm:"type:REAL;"`
	Ampere    float32        `json:"ampere" gorm:"type:REAL"`
	Watt      float32        `json:"watt" gorm:"type:REAL;"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
