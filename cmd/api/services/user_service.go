package services

import (
	"antara-api/cmd/api/requests"
	"antara-api/common"
	"antara-api/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (userService *UserService) CreateUser(signUpRequest *requests.SignUpRequest) (*models.UserModel, error) {
	hashedPassword, err := common.HashPassword(signUpRequest.Password)
	if err != nil {
		return nil, errors.New("register failed")
	}

	createdUser := models.UserModel{
		Name:     &signUpRequest.Name,
		Email:    signUpRequest.Email,
		Password: hashedPassword,
	}
	result := userService.db.Create(&createdUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &createdUser, nil
}

func (userService *UserService) GetByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	result := userService.db.Preload("Roles").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService *UserService) GetById(id uint) (*models.UserModel, error) {
	var user models.UserModel
	result := userService.db.Preload("Roles").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userService *UserService) ChangeUserPassword(user *models.UserModel, newPassword string) error {
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		fmt.Println(err)
		return errors.New("password change failed")
	}

	result := userService.db.Model(user).Update("Password", hashedPassword)
	if result.Error != nil {
		return errors.New("password change failed")
	}

	return nil
}
