package service

import (
	r "book-service/config"
	"book-service/model"
	"book-service/repository"
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type BookService interface {
	CreateBook(req *model.CreateBookRequest, sellerID uint) (*model.BookResponse, error)
	GetAllBooks(category string) ([]model.BookResponse, error)
	GetBookByID(id uint) (*model.BookResponse, error)
	GetBooksBySellerID(sellerID uint) ([]model.BookResponse, error)
	UpdateBook(id uint, req *model.UpdateBookRequest, sellerID uint) (*model.BookResponse, error)
	DeleteBook(id uint, sellerID uint) error
	DeductStock(id uint, amount int) (*model.BookResponse, error)
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

	// save to redis
	// Save to Redis (HSET)
	err = r.Client.HSet(
		context.Background(),
		fmt.Sprintf("book:%d", book.ID),
		map[string]interface{}{
			"id":          book.ID,
			"name":        book.Name,
			"description": book.Description,
			"author":      book.Author,
			"stock":       book.Stock,
			"costs":       book.Costs,
			"category":    book.Category,
			"seller_id":   book.SellerID,
		},
	).Err()

	if err != nil {
		fmt.Println("⚠ Failed to cache book:", err)
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
	ctx := context.Background()
	key := fmt.Sprintf("book:%d", id)

	// Try to get from Redis first
	data, err := r.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// If Redis returned empty, fallback to DB
	if len(data) == 0 {
		book, err := s.bookRepo.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("book not found")
			}
			return nil, err
		}
		// Save to Redis for next time
		r.Client.HSet(ctx, key, map[string]interface{}{
			"id":          book.ID,
			"seller_id":   book.SellerID,
			"name":        book.Name,
			"description": book.Description,
			"author":      book.Author,
			"stock":       book.Stock,
			"costs":       book.Costs,
			"category":    book.Category,
		})

		response := book.ToResponse()
		return &response, nil
	}

	// Map Redis hash back to BookResponse
	stock, _ := strconv.Atoi(data["stock"])
	costs, _ := strconv.Atoi(data["costs"])
	sellerID, _ := strconv.Atoi(data["seller_id"])
	bookID, _ := strconv.Atoi(data["id"])

	response := &model.BookResponse{
		ID:          uint(bookID),
		SellerID:    uint(sellerID),
		Name:        data["name"],
		Description: data["description"],
		Author:      data["author"],
		Stock:       stock,
		Costs:       float64(costs),
		Category:    data["category"],
	}

	return response, nil
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

func (s *bookService) DeductStock(id uint, amount int) (*model.BookResponse, error) {
	if amount <= 0 {
		return nil, errors.New("deduction amount must be greater than 0")
	}

	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if book.Stock < amount {
		return nil, errors.New("insufficient stock")
	}

	err = s.bookRepo.DeductStock(id, amount)
	if err != nil {
		return nil, err
	}

	updatedBook, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := updatedBook.ToResponse()
	return &response, nil

}
