package handler

import (
	"auth-service/dto"
	"auth-service/helpers"
	"auth-service/models"
	"auth-service/service"
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	// Validate the input
	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Validation failed: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	createdUser, err := h.Service.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create user: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "User registered successfully",
		User:    createdUser,
	})
}
func (h *AuthHandler) Login(c echo.Context) error {
	var input dto.LoginRequest

	// Parse JSON input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := h.Service.GetUserByEmail(input.Email)
	if err != nil {

		if err.Error() == "email not found" {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Email not found",
				Code:    http.StatusBadRequest,
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to retrieve user: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		// If password does not match, return unauthorized
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid email or password",
			Code:    http.StatusUnauthorized,
		})
	}
	// Generate JWT token
	token, err := helpers.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate token: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := h.Service.GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
			Code:    http.StatusNotFound,
		})
	}
	return c.JSON(http.StatusOK, dto.GetUserByIDResponse{
		Message: "User retrieved successfully",
		User:    user,
	})
}

func (h *AuthHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	user.ID = uint(id)
	updatedUser, err := h.Service.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, updatedUser)
}
