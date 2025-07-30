package dto

import "auth-service/models"

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RegisterResponse struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

type GetUserByIDResponse struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

type UpdateBalanceResponse struct {
	Message string  `json:"message"`
	ID      uint    `json:"id"`
	Balance float64 `json:"balance"`
}

type VerificationResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}
