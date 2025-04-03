package services

import (
	"antara-api/internal/models"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type VerificationTokenService struct {
	db *gorm.DB
}

func NewVerificationTokenService(db *gorm.DB) *VerificationTokenService {
	return &VerificationTokenService{db: db}
}

func (VerificationTokenService *VerificationTokenService) getToken() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(99999-10000+1) + 10000
}

func (VerificationTokenService *VerificationTokenService) GenerateResetPasswordToken(user models.UserModel) (*models.VerificationTokenModel, error) {
	expirationSecondsStr := os.Getenv("VERIFY_EMAIL_TOKEN_EXPIRATION_SECONDS")
	expirationSecondsInt, err := strconv.Atoi(expirationSecondsStr)

	if err != nil {
		return nil, errors.New("invalid expiration seconds")
	}

	tokenCreated := models.VerificationTokenModel{
		TargetId:   user.ID,
		Identifier: "reset_password",
		Token:      strconv.Itoa(VerificationTokenService.getToken()),
		Used:       false,
		ExpiresAt:  time.Now().Add(time.Second * time.Duration(expirationSecondsInt)),
	}

	result := VerificationTokenService.db.Create(&tokenCreated)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tokenCreated, nil
}

func (VerificationTokenService *VerificationTokenService) ValidateResetPasswordToken(user models.UserModel, token string) (*models.VerificationTokenModel, error) {
	var existingToken models.VerificationTokenModel

	result := VerificationTokenService.db.Where(&models.VerificationTokenModel{
		TargetId:   user.ID,
		Identifier: "reset_password",
		Token:      token,
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

func (VerificationTokenService *VerificationTokenService) InvalidateToken(userId uint, VerificationToken models.VerificationTokenModel) {
	VerificationTokenService.db.Model(&models.VerificationTokenModel{}).Where("target_id = ? AND token = ?", userId, VerificationToken.Token).Update("used", true)
}
