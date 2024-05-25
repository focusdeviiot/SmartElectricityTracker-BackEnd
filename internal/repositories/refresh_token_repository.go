package repositories

import (
	"errors"
	"smart_electricity_tracker_backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(refreshToken *models.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

func (r *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		r.DeleteRefreshToken(&refreshToken)
		return nil, errors.New("refresh token expired")
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(refreshToken *models.RefreshToken) error {
	return r.db.Delete(refreshToken).Error
}
