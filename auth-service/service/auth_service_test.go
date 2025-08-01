package service

import (
	"auth-service/dto"
	"auth-service/models"
	"auth-service/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserByID(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	expectedUser := models.User{
		ID:       1,
		Fullname: "test user",
		Email:    "test@mail.com",
		Role:     "seller",
		Address:  "123 Street",
	}

	mockRepo.On("GetUserByID", uint(1)).Return(expectedUser, nil)

	users, err := svc.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, users)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	expectedUser := models.User{
		ID:       2,
		Fullname: "test user",
		Email:    "test@mail.com",
		Role:     "seller",
		Address:  "123 Street",
	}
	mockRepo.On("GetUserByEmail", "test@mail.com").Return(expectedUser, nil)
	user, err := svc.GetUserByEmail("test@mail.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		FullName: "John Doe",
		Email:    "john@example.com",
		Password: "securepass",
		Address:  "Somewhere",
		Role:     "buyer",
	}

	expectedUser := models.User{
		ID:       1,
		Fullname: "John Doe",
		Email:    "john@example.com",
		Address:  "Somewhere",
		Role:     "buyer",
	}

	mockRepo.On("CreateUser", mock.MatchedBy(func(user models.User) bool {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		return user.Email == req.Email && err == nil
	})).Return(expectedUser, nil)

	result, err := svc.CreateUser(req)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	inputUser := models.User{
		ID:       1,
		Fullname: "Jane Doe",
		Email:    "jane@example.com",
		Role:     "buyer",
		Address:  "Old Address",
	}

	expectedUser := inputUser
	expectedUser.Address = "New Address"

	mockRepo.On("UpdateUser", inputUser).Return(expectedUser, nil)

	result, err := svc.UpdateUser(inputUser)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestVerifyUser_AlreadyVerified(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo)

	existingUser := models.User{
		Email:      "user@example.com",
		IsVerified: true,
	}

	mockRepo.On("VerifyUser", existingUser.Email).Return(existingUser, nil)

	user, err := svc.VerifyUser(existingUser.Email)
	assert.NoError(t, err)
	assert.Equal(t, true, user.IsVerified)
	mockRepo.AssertNotCalled(t, "UpdateUser")
}
