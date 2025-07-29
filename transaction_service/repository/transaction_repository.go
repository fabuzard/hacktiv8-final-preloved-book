package repository

import (
	"main/model"
	"main/utils"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(user_id int, t model.Transaction) (model.Transaction, error)
	GetTransaction(user_id int) ([]model.Transaction, error)
	UpdateTransactionStatus(transaction_id int) (model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(user_id int, t model.Transaction) (model.Transaction, error) {
	t.User_ID = user_id
	t.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	t.Expiration_Date = time.Now().Add(4 * time.Hour).Format("2006-01-02 15:04:05")
	t.Status = "pending"
	if err := r.db.Create(&t).Error; err != nil {
		return model.Transaction{}, err
	}
	return t, nil
}

func (r *transactionRepository) GetTransaction(user_id int) ([]model.Transaction, error) {
	var trans []model.Transaction
	err := r.db.Where("user_id = ?", user_id).Find(&trans).Error

	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	return trans, nil
}

func (r *transactionRepository) UpdateTransactionStatus(transaction_id int) (model.Transaction, error) {
	var t model.Transaction
	err := r.db.Model(&t).Where("transaction_id = ?", transaction_id).Update("status", "success").Error
	if err != nil {
		return model.Transaction{}, utils.ErrBadReq
	}

	if err := r.db.First(&t, transaction_id).Error; err != nil {
		return model.Transaction{}, utils.ErrUserNotFound
	}

	return t, nil
}
