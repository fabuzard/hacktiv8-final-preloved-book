package service

import (
	"book-service/model"
	"book-service/repository"
	"errors"
	"gorm.io/gorm"
)

type BookService interface {
	CreateBook(req *model.CreateBookRequest, sellerID uint) (*model.BookResponse, error)
	GetAllBooks(category string) ([]model.BookResponse, error)
	GetBookByID(id uint) (*model.BookResponse, error)
	GetBooksBySellerID(sellerID uint) ([]model.BookResponse, error)
	UpdateBook(id uint, req *model.UpdateBookRequest, sellerID uint) (*model.BookResponse, error)
	DeleteBook(id uint, sellerID uint) error
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (s *bookService) CreateBook(req *model.CreateBookRequest, sellerID uint) (*model.BookResponse, error) {
	book := &model.Book{
		SellerID:    sellerID,
		Name:        req.Name,
		Description: req.Description,
		Author:      req.Author,
		Stock:       req.Stock,
		Costs:       req.Costs,
		Category:    req.Category,
	}

	err := s.bookRepo.Create(book)
	if err != nil {
		return nil, err
	}

	response := book.ToResponse()
	return &response, nil
}

func (s *bookService) GetAllBooks(category string) ([]model.BookResponse, error) {
	books, err := s.bookRepo.GetAll(category)
	if err != nil {
		return nil, err
	}

	var responses []model.BookResponse
	for _, book := range books {
		responses = append(responses, book.ToResponse())
	}

	return responses, nil
}

func (s *bookService) GetBookByID(id uint) (*model.BookResponse, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	response := book.ToResponse()
	return &response, nil
}

func (s *bookService) GetBooksBySellerID(sellerID uint) ([]model.BookResponse, error) {
	books, err := s.bookRepo.GetBySellerID(sellerID)
	if err != nil {
		return nil, err
	}

	var responses []model.BookResponse
	for _, book := range books {
		responses = append(responses, book.ToResponse())
	}

	return responses, nil
}

func (s *bookService) UpdateBook(id uint, req *model.UpdateBookRequest, sellerID uint) (*model.BookResponse, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if book.SellerID != sellerID {
		return nil, errors.New("unauthorized: you can only update your own books")
	}

	if req.Name != nil {
		book.Name = *req.Name
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Stock != nil {
		book.Stock = *req.Stock
	}
	if req.Costs != nil {
		book.Costs = *req.Costs
	}
	if req.Category != nil {
		book.Category = *req.Category
	}

	err = s.bookRepo.Update(book)
	if err != nil {
		return nil, err
	}

	response := book.ToResponse()
	return &response, nil
}

func (s *bookService) DeleteBook(id uint, sellerID uint) error {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("book not found")
		}
		return err
	}

	if book.SellerID != sellerID {
		return errors.New("unauthorized: you can only delete your own books")
	}

	return s.bookRepo.Delete(id, sellerID)
}