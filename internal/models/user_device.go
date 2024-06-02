package models

import (
	"time"

	"gorm.io/gorm"
)

type UserDevice struct {
	gorm.Model
	ID            string         `json:"id" gorm:"type:varchar(255);default:gen_random_uuid();primary_key;index:;"`
	UserID        string         `json:"user_id" gorm:"type:varchar(255);index:;not null;"`
	User          User           `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeviceID      string         `json:"device_id" gorm:"type:varchar(255);index:;not null;"`
	Device        DeviceMaster   `json:"device" gorm:"foreignKey:DeviceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedBy     string         `json:"created_by" gorm:"type:varchar(255);default:'SYSTEM';index:;not null;"`
	CreatedByName string         `json:"created_by_name" gorm:"type:varchar(255);default:'SYSTEM';index:;not null;"`
	UpdateBy      string         `json:"update_by" gorm:"type:varchar(255);default:'SYSTEM';index:;not null;"`
	UpdateByName  string         `json:"update_by_name" gorm:"type:varchar(255);default:'SYSTEM';index:;not null;"`
	CreatedAt     time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
