package handler

import (
	"auth-service/dto"
	"auth-service/helpers"
	"auth-service/models"
	"auth-service/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: s,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var user dto.RegisterRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	createdUser, err := h.Service.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdUser)
}
func (h *AuthHandler) Login(c echo.Context) error {
	var input dto.LoginRequest

	// Parse JSON input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	user, err := h.Service.GetUserByEmail(input.Email)
	if err != nil {

		if err.Error() == "email not found" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email not registered"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Server error"})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {

		fmt.Println("stored hash: ", user.Password)
		fmt.Println("input password: ", input.Password)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}
	// Generate JWT token
	token, err := helpers.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		Message: "Login successful",
		Token:   token,
	})
}

func (h *AuthHandler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	user, err := h.Service.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	user.ID = uint(id)
	updatedUser, err := h.Service.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updatedUser)
}
