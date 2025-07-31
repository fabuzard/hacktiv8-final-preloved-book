package service

import (
	"auth-service/models"
	"auth-service/repository"
	"testing"

	"github.com/stretchr/testify/assert"
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
