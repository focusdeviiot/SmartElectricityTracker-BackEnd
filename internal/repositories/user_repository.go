package repositories

import (
	"smart_electricity_tracker_backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUserId(userId uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteUser(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *UserRepository) FindUsersCountDevice() ([]models.UserCountDevice, error) {
	var users []models.UserCountDevice
	if err := r.db.Table("users").Select("users.id as user_id, users.username, users.name, users.role, count(user_devices.id) as device_count").
		Joins("left join user_devices on user_devices.user_id = users.id").
		Group("users.id").
		Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindUserCountDeviceById(userId uuid.UUID) (*models.UserDevice, error) {
	var user models.UserDevice
	if err := r.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
