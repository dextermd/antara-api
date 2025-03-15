package services

import (
	"antara-api/internal/models"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

type AppTokenService struct {
	db *gorm.DB
}

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

func (appTokenService *AppTokenService) getToken() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(99999-10000+1) + 10000
}

func (appTokenService *AppTokenService) GenerateResetPasswordToken(user models.UserModel) (*models.AppTokenModel, error) {
	tokenCreated := models.AppTokenModel{
		TargetId:  user.Id,
		Type:      "reset_password",
		Token:     strconv.Itoa(appTokenService.getToken()),
		Used:      false,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	result := appTokenService.db.Create(&tokenCreated)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tokenCreated, nil
}

func (appTokenService *AppTokenService) ValidateResetPasswordToken(user models.UserModel, token string) (*models.AppTokenModel, error) {
	var existingToken models.AppTokenModel

	result := appTokenService.db.Where(&models.AppTokenModel{
		TargetId: user.Id,
		Type:     "reset_password",
		Token:    token,
	}).First(&existingToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid password reset token")
		}
		return nil, result.Error
	}

	if existingToken.Used {
		return nil, errors.New("password reset token already used")
	}

	if existingToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("password reset token expired, please re-initiate forgot password")
	}

	return &existingToken, nil
}

func (appTokenService *AppTokenService) InvalidateToken(userId uint, appToken models.AppTokenModel) {
	appTokenService.db.Model(&models.AppTokenModel{}).Where("target_id = ? AND token = ?", userId, appToken.Token).Update("used", true)
}
