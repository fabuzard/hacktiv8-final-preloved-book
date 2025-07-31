package dto

import "main/model"

type GetBookByIDResponse struct {
	Message string       `json:"message"`
	Data    BookResponse `json:"data"`
}

type BookResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Cost  float64 `json:"costs"`
}

type GetUserByIDResponse struct {
	Message string     `json:"message"`
	User    model.User `json:"user"`
}
