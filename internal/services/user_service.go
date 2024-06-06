package services

import (
	"smart_electricity_tracker_backend/internal/config"
	"smart_electricity_tracker_backend/internal/models"
	"smart_electricity_tracker_backend/internal/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo                 *repositories.UserRepository
	refreshTokenRepo     *repositories.RefreshTokenRepository
	jwtSecret            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	cfg                  *config.Config
}

func NewUserService(repo *repositories.UserRepository, refreshTokenRepo *repositories.RefreshTokenRepository, jwtSecret string, accessTokenDuration time.Duration, refreshTokenDuration time.Duration, cfg *config.Config) *UserService {
	return &UserService{repo: repo, refreshTokenRepo: refreshTokenRepo, jwtSecret: jwtSecret, accessTokenDuration: accessTokenDuration, refreshTokenDuration: refreshTokenDuration, cfg: cfg}
}

func (s *UserService) GenerateTokens(user *models.User) (string, string, error) {
	accessToken, err := s.createAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.createRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) createAccessToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
		Exp:      time.Now().Add(s.accessTokenDuration),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) createRefreshToken(user *models.User) (string, error) {
	token := uuid.New().String()
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
	}

	if err := s.refreshTokenRepo.CreateRefreshToken(refreshToken); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Authenticate(username, password string) (string, string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", err
	}

	return s.GenerateTokens(user)
}

func (s *UserService) RefreshToken(refreshTokenString string) (string, string, error) {
	refreshToken, err := s.refreshTokenRepo.FindByToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.FindByUserId(refreshToken.UserID)
	if err != nil {
		return "", "", err
	}

	if err := s.refreshTokenRepo.DeleteRefreshToken(refreshToken); err != nil {
		return "", "", err
	}

	return s.GenerateTokens(user)
}

func (s *UserService) CreateUser(username, password, role, name string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	roleCon, err := models.StringToRole(role)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     roleCon,
		Name:     name,
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) Logout(refreshTokenString string) error {
	refreshToken, err := s.refreshTokenRepo.FindByToken(refreshTokenString)
	if err != nil {
		return err
	}

	return s.refreshTokenRepo.DeleteRefreshToken(refreshToken)
}

func (s *UserService) GetUsers() ([]models.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUserById(id string) (*models.GetUserRes, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	userGet, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	userRes := &models.GetUserRes{
		UserID:   userGet.ID.String(),
		Username: userGet.Username,
		Name:     userGet.Name,
		Role:     userGet.Role,
	}

	return userRes, nil
}

func (s *UserService) UpdateUser(id string, password, role, name string) error {
	userId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	user, err := s.repo.FindByUserId(userId)
	if err != nil {
		return err
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	if role != "" {
		roleCon, err := models.StringToRole(role)
		if err != nil {
			return err
		}
		user.Role = roleCon
	}

	if name != "" {
		user.Name = name
	}

	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id string) error {
	userId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	user, err := s.repo.FindByUserId(userId)
	if err != nil {
		return err
	}

	return s.repo.DeleteUser(user)
}

func (s *UserService) GetUserByUsername(username string) (*models.GetUserRes, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	userRes := &models.GetUserRes{
		UserID:   user.ID.String(),
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}

	return userRes, nil
}

func (s *UserService) GetAllUsersCountDevice(req *models.SearchUserCountDeviceListReq) ([]models.UserCountDeviceRes, *models.Pageable, error) {
	return s.repo.FindAllUsersCountDevice(req)
}

func (s *UserService) GetUserDeviceById(id *string) ([]models.UserDeviceFromDB, error) {
	userId, err := uuid.Parse(*id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindUserDeviceById(userId)
}

func (s *UserService) UpdateUserDevice(req *models.UpdateUserDeviceReq) error {
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		return err
	}

	return s.repo.UpdateUserDevice(userId, req.DeviceID)
}
