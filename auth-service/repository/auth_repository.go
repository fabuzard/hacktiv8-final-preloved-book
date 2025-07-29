package repository

import (
	"auth-service/models"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
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

func (r *authRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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
