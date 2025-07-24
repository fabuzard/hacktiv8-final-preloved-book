package service

import (
	"auth-service/dto"
	"auth-service/models"
	"auth-service/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (models.User, error)
	CreateUser(user dto.RegisterRequest) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) GetUserByEmail(email string) (models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *authService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *authService) GetUserByID(id uint) (models.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (s *authService) CreateUser(user dto.RegisterRequest) (models.User, error) {

	// Convert dto.RegisterRequest to models.User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	InputUser := models.User{
		Fullname: user.FullName,
		Email:    user.Email,
		Password: string(hashedPassword),
		Address:  user.Address,
		Role:     user.Role,
	}
	createdUser, err := s.repo.CreateUser(InputUser)
	if err != nil {
		return models.User{}, err
	}
	return createdUser, nil
}
func (s *authService) UpdateUser(user models.User) (models.User, error) {
	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}
