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
