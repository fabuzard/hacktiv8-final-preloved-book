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

// Other methods are not used in this test so we just make them panic
func (m *MockAuthService) CreateUser(dto.RegisterRequest) (models.User, error) {
	panic("not implemented")
}
func (m *MockAuthService) UpdateUser(user models.User) (models.User, error) {
	if user.ID == 99 {
		return models.User{}, errors.New("user not found")
	}
	user.Fullname = "Updated Name"
	return user, nil
}
func (m *MockAuthService) GetUserByEmail(email string) (models.User, error) {
	panic("not implemented")
}
func (m *MockAuthService) DeleteUser(id uint) error {
	panic("not implemented")
}
func (m *MockAuthService) VerifyUser(email string) (models.User, error) {
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
