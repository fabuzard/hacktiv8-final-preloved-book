package repository

import (
	"auth-service/models"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByID(id uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteInactiveUsersOver30Days() error

	VerifyUser(email string) (models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) DeleteInactiveUsersOver30Days() error {
	return r.db.
		Where("is_verified = ? AND created_at <= NOW() - INTERVAL '30 days'", false).
		Delete(&models.User{}).Error
}

func (r *authRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("email not found")
		}

		return models.User{}, err
	}
	return user, nil
}

func (r *authRepository) GetUserByID(id uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *authRepository) CreateUser(user models.User) (models.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *authRepository) UpdateUser(user models.User) (models.User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *authRepository) VerifyUser(email string) (models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("email not found")
		}
		return models.User{}, err
	}

	user.IsVerified = true
	if err := r.db.Save(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
