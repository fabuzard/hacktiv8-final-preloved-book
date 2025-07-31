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

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Fullname   string    `gorm:"type:varchar(100);not null"`
	Email      string    `gorm:"type:varchar(100);unique;not null"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	Address    string    `gorm:"type:text"`
	Role       string    `gorm:"type:varchar(10);not null"`
	Balance    float64   `gorm:"type:decimal(12,2);default:0.00"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	IsVerified bool      `gorm:"default:false"`
}
