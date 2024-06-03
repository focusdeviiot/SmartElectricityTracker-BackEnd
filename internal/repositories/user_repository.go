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

func buildUserCountDeviceQuery(db *gorm.DB, req *models.SearchUserCountDeviceListReq) *gorm.DB {
	query := db.Table("users as u")
	query = query.Joins("left join user_devices ud2 on u.id = ud2.user_id")
	query = query.Joins("left join (select ud1.user_id, count(ud1.id) as count from user_devices ud1 group by ud1.user_id) ud1 on ud1.user_id = u.id")
	query = query.Where("u.deleted_at IS NULL")

	if req.Username != "" {
		query = query.Where("u.username ILIKE ?", "%"+req.Username+"%")
	}

	if req.Name != "" {
		query = query.Where("u.name ILIKE ?", "%"+req.Name+"%")
	}

	if req.Role != "*" {
		query = query.Where("u.role = ?", req.Role)
	}

	if req.DeviceId != "*" {
		query = query.Where("ud2.device_id = ?", req.DeviceId)
	}

	return query
}

func (r *UserRepository) FindUsersCountDevice(req *models.SearchUserCountDeviceListReq) ([]models.UserCountDeviceRes, *models.Pageable, error) {
	var users []models.UserCountDeviceRes

	pageable := &models.Pageable{}
	baseQuery := buildUserCountDeviceQuery(r.db, req)

	// Count distinct users
	var totalElements int64
	countQuery := buildUserCountDeviceQuery(r.db, req)
	err := countQuery.Select("count(distinct u.id)").Scan(&totalElements).Error
	if err != nil {
		return nil, pageable, err
	}

	query := baseQuery.Select("distinct u.id as user_id, u.username, u.name, u.role, ud1.count as device_count")
	query = query.Order("u.username ASC, u.name ASC, u.role ASC")

	pageable.PageNumber = int(req.Pageable.PageNumber)
	pageable.PageSize = int(req.Pageable.PageSize)
	pageable.TotalElements = totalElements

	if req.Pageable.PageSize > 0 {
		pageable.TotalPages = int((totalElements + int64(req.Pageable.PageSize) - 1) / int64(req.Pageable.PageSize))
		query = query.Offset(int((req.Pageable.PageNumber - 1) * req.Pageable.PageSize)).Limit(int(req.Pageable.PageSize))
	} else {
		pageable.TotalPages = 1
	}

	if err := query.Scan(&users).Error; err != nil {
		return nil, nil, err
	}
	return users, pageable, nil
}

func (r *UserRepository) FindUserCountDeviceById(userId uuid.UUID) (*models.UserDevice, error) {
	var user *models.UserDevice
	if err := r.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
