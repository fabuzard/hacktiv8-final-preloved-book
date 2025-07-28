package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Fullname  string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Address   string    `gorm:"type:text"`
	Role      string    `gorm:"type:varchar(10);not null"`
	Balance   float64   `gorm:"type:decimal(12,2);default:0.00"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
