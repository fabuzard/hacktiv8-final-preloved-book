// handler/auth_handler_test.go
package handler

import (
	"auth-service/dto"
	"auth-service/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Mock service
type MockAuthService struct{}

func (m *MockAuthService) GetUserByID(id uint) (models.User, error) {
	return models.User{
		ID:       id,
		Fullname: "Fahreza",
		Email:    "reza@mail.com",
		Role:     "customer",
	}, nil
}

func (m *MockAuthService) CreateUser(dto.RegisterRequest) (models.User, error) {
	return models.User{}, nil
}
func (m *MockAuthService) UpdateUser(user models.User) (models.User, error) {
	if user.ID == 99 {
		return models.User{}, errors.New("user not found")
	}
	user.Fullname = "Updated Name"
	return user, nil
}
func (m *MockAuthService) GetUserByEmail(email string) (models.User, error) {
	if email == "notfound@mail.com" {
		return models.User{}, errors.New("user not found")
	}
	return models.User{
		ID:       2,
		Fullname: "Email User",
		Email:    email,
		Role:     "buyer",
	}, nil
}

func (m *MockAuthService) DeleteUser(id uint) error {
	if id == 999 {
		return errors.New("user not found")
	}
	return nil
}

func (m *MockAuthService) VerifyUser(email string) (models.User, error) {
	panic("not implemented")
}
func (m *MockAuthService) DeleteInactiveUsersOver30Days() error {
	panic("not implemented")
}

func TestGetUserByID(t *testing.T) {
	e := echo.New()

	// Create mock service and inject into handler
	mockService := &MockAuthService{}
	h := &AuthHandler{Service: mockService}

	// Prepare the request
	userID := uint(10)
	req := httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(int(userID)), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(userID)))

	// Call the handler
	if assert.NoError(t, h.GetUserByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Fahreza")
		assert.Contains(t, rec.Body.String(), "reza@mail.com")
	}
}

func TestUpdateUser(t *testing.T) {
	e := echo.New()

	mockService := &MockAuthService{}
	handler := &AuthHandler{Service: mockService}

	// Sample user payload
	userPayload := models.User{
		Fullname: "Old Name",
		Email:    "test@example.com",
		Address:  "123 Test St",
	}

	jsonBody, _ := json.Marshal(userPayload)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Set context
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(1))

	// Act
	if assert.NoError(t, handler.UpdateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Updated Name")
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	mock := &MockAuthService{}
	email := "found@mail.com"

	user, err := mock.GetUserByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "Email User", user.Fullname)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	mock := &MockAuthService{}
	email := "notfound@mail.com"

	user, err := mock.GetUserByEmail(email)

	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
	assert.Equal(t, uint(0), user.ID)
}

func TestDeleteUser_Success(t *testing.T) {
	mock := &MockAuthService{}
	err := mock.DeleteUser(1)

	assert.NoError(t, err)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mock := &MockAuthService{}
	err := mock.DeleteUser(999)

	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
}
