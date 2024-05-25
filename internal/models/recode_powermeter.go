package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecodePowermeter struct {
	gorm.Model
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key;index:;"`
	DeviceID  string         `json:"device_id" gorm:"type:varchar(255);index:;not null;"`
	Device    DeviceMaster   `json:"device" gorm:"foreignKey:DeviceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Volt      int            `json:"volt" gorm:"type:integer;index:;not null;"`
	Ampere    int            `json:"ampere" gorm:"type:integer;index:;not null;"`
	Watt      int            `json:"watt" gorm:"type:integer;index:;not null;"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
