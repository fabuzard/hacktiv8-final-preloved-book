package model

import (
	"time"
	"gorm.io/gorm"
)

type Book struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	SellerID    uint           `json:"seller_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null;size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Author      string         `json:"author" gorm:"size:100"`
	Stock       int            `json:"stock" gorm:"default:0"`
	Costs       float64        `json:"costs" gorm:"not null;type:decimal(10,2)"`
	Category    string         `json:"category" gorm:"size:100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateBookRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	Stock       int     `json:"stock" validate:"min=0"`
	Costs       float64 `json:"costs" validate:"required,min=0"`
	Category    string  `json:"category"`
}

type UpdateBookRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Author      *string  `json:"author,omitempty"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	Costs       *float64 `json:"costs,omitempty" validate:"omitempty,min=0"`
	Category    *string  `json:"category,omitempty"`
}

type BookResponse struct {
	ID          uint      `json:"id"`
	SellerID    uint      `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Stock       int       `json:"stock"`
	Costs       float64   `json:"costs"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

func (b *Book) ToResponse() BookResponse {
	return BookResponse{
		ID:          b.ID,
		SellerID:    b.SellerID,
		Name:        b.Name,
		Description: b.Description,
		Author:      b.Author,
		Stock:       b.Stock,
		Costs:       b.Costs,
		Category:    b.Category,
		CreatedAt:   b.CreatedAt,
	}
}