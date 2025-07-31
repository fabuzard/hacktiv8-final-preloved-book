package service

import (
	"main/model"
	"main/repository"
)

type TransactionService interface {
	CreateTransaction(user_id int, t model.Transaction) (model.Transaction, error)
	GetTransaction(user_id int) ([]model.Transaction, error)
	UpdateTransactionStatus(transaction_id int) (model.Transaction, error)
	GetTransactionByID(transaction_id int) (model.Transaction, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(user_id int, t model.Transaction) (model.Transaction, error) {
	trans, err := s.repo.CreateTransaction(user_id, t)
	if err != nil {
		return model.Transaction{}, err
	}
	return trans, nil
}

func (s *transactionService) GetTransaction(user_id int) ([]model.Transaction, error) {
	trans, err := s.repo.GetTransaction(user_id)
	if err != nil {
		return nil, err
	}
	return trans, nil
}

func (s *transactionService) UpdateTransactionStatus(transaction_id int) (model.Transaction, error) {
	trans, err := s.repo.UpdateTransactionStatus(transaction_id)
	if err != nil {
		return model.Transaction{}, err
	}
	return trans, nil
}

func (s *transactionService) GetTransactionByID(transaction_id int) (model.Transaction, error) {
	trans, err := s.repo.GetTransactionByID(transaction_id)
	if err != nil {
		return model.Transaction{}, err
	}
	return trans, nil
}
