package handler

import (
	"auth-service/dto"
	"auth-service/helpers"
	"auth-service/models"
	"auth-service/service"
	"log"
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

	// Generate email verification token
	emailToken, err := helpers.GenerateEmailToken(createdUser.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate email verification token: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	go func() {
		err := helpers.SendVerificationEmail(user.Email, emailToken)
		if err != nil {
			log.Println("failed to send verification email:", err)
		}
	}()

	return c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "User registered successfully, Please check your email for verification",
		User:    createdUser,
	})
}

func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	var req dto.ResendVerificationEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := h.Service.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
			Code:    http.StatusNotFound,
		})
	}

	if user.IsVerified {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "User is already verified",
			Code:    http.StatusBadRequest,
		})
	}

	emailToken, err := helpers.GenerateEmailToken(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate email verification token: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	go func() {
		err := helpers.SendVerificationEmail(user.Email, emailToken)
		if err != nil {
			log.Println("failed to send verification email:", err)
		}
	}()

	return c.JSON(http.StatusOK, dto.VerificationResponse{
		Message: "Verification email sent successfully",
		Email:   req.Email,
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

	if !user.IsVerified {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "User is not verified, please check your email or send another request for verification",
			Code:    http.StatusUnauthorized,
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
	token, err := helpers.GenerateJWT(user)
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

func (h *AuthHandler) UpdateBalance(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	// Bind the request body to UpdateBalanceRequest
	var balanceRequest dto.UpdateBalanceRequest
	if err := c.Bind(&balanceRequest); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
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

	if balanceRequest.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Amount must be greater than zero",
			Code:    http.StatusBadRequest,
		})
	}

	user.Balance += balanceRequest.Amount
	updatedUser, err := h.Service.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update user balance: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, dto.UpdateBalanceResponse{
		Message: "User balance updated successfully",
		ID:      updatedUser.ID,
		Balance: updatedUser.Balance,
	})

}

func (h *AuthHandler) VerifyUser(c echo.Context) error {
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Token is required",
			Code:    http.StatusBadRequest,
		})
	}

	email, err := helpers.ParseAndValidateEmailToken(tokenString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Token verification failed: " + err.Error(),
			Code:    http.StatusUnauthorized,
		})
	}

	user, err := h.Service.VerifyUser(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to verify user: " + err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, dto.RegisterResponse{
		Message: "User verified successfully",
		User:    user,
	})
}
