package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	Transaction_ID  uint           `gorm:"primaryKey;autoincrement" json:"transaction_id"`
	Amount          float64        `gorm:"not null" json:"amount"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	User_ID         int            `gorm:"not null" json:"user_id"`
	Status          string         `gorm:"not null" json:"status"`
	Book_ID         int            `gorm:"not null" json:"book_id"`
	Expiration_Date time.Time      `gorm:"not null" json:"expiration_date"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}
