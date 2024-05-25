package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key;index:;"`
	Username  string         `json:"username" gorm:"unique;index:;not null;type:varchar(50)"`
	Password  string         `json:"password" gorm:"not null;type:varchar(255)"`
	Name      string         `json:"name" gorm:"not null;type:varchar(50)"`
	Role      Role           `json:"role" gorm:"type:user_role;default:'USER'"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
