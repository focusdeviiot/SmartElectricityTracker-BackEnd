package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceMaster struct {
	ID        string         `json:"id" gorm:"type:varchar(50);primary_key;index:;"`
	Name      string         `json:"name" gorm:"unique;type:varchar(255);index:;not null;"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
