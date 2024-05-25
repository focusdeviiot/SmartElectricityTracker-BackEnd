package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key;index:;"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;index:;not null;"`
	Token     string    `gorm:"unique; not null; type:varchar(255); index:;"`
	ExpiresAt time.Time
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
