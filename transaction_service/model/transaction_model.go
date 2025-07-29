package model

type Transaction struct {
	Transaction_ID    uint    `gorm:"primaryKey" json:"transaction_id"`
	Name              string  `gorm:"not null" json:"name"`
	Amount            float64 `gorm:"not null" json:"amount"`
	CreatedAt         string  `gorm:"not null" json:"created_at"`
	DeletedAt         string  `gorm:"not null" json:"deleted_at"`
	User_ID           int     `gorm:"not null" json:"user_id"`
	Status            string  `gorm:"not null" json:"status"`
	Book_ID           int     `gorm:"not null" json:"book_id"`
	Expipiration_Date string  `gorm:"not null" json:"expiration_date"`
}
