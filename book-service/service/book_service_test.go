package service

import (
	"book-service/model"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) GetAll(category string) ([]model.Book, error) {
	args := m.Called(category)
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookRepository) GetByID(id uint) (*model.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *MockBookRepository) GetBySellerID(sellerID uint) ([]model.Book, error) {
	args := m.Called(sellerID)
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookRepository) Update(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(id uint, sellerID uint) error {
	args := m.Called(id, sellerID)
	return args.Error(0)
}

func (m *MockBookRepository) DeductStock(id uint, amount int) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func TestCreateBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	req := &model.CreateBookRequest{
		Name:        "Test Book",
		Description: "Test Description",
		Author:      "Test Author",
		Stock:       10,
		Costs:       29.99,
		Category:    "Fiction",
	}

	expectedBook := &model.Book{
		ID:          1,
		SellerID:    1,
		Name:        req.Name,
		Description: req.Description,
		Author:      req.Author,
		Stock:       req.Stock,
		Costs:       req.Costs,
		Category:    req.Category,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Book")).Return(nil).Run(func(args mock.Arguments) {
		book := args.Get(0).(*model.Book)
		book.ID = 1
	})

	result, err := service.CreateBook(req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBook.Name, result.Name)
	assert.Equal(t, expectedBook.SellerID, result.SellerID)
	mockRepo.AssertExpectations(t)
}

func TestCreateBook_RepositoryError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	req := &model.CreateBookRequest{
		Name:        "Test Book",
		Description: "Test Description",
		Author:      "Test Author",
		Stock:       10,
		Costs:       29.99,
		Category:    "Fiction",
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Book")).Return(errors.New("database error"))

	result, err := service.CreateBook(req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetAllBooks_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	expectedBooks := []model.Book{
		{ID: 1, Name: "Book 1", SellerID: 1, Costs: 29.99},
		{ID: 2, Name: "Book 2", SellerID: 2, Costs: 39.99},
	}

	mockRepo.On("GetAll", "fiction").Return(expectedBooks, nil)

	result, err := service.GetAllBooks("fiction")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedBooks[0].Name, result[0].Name)
	assert.Equal(t, expectedBooks[1].Name, result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestGetAllBooks_RepositoryError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("GetAll", "").Return([]model.Book{}, errors.New("database error"))

	result, err := service.GetAllBooks("")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetBookByID_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	expectedBook := &model.Book{
		ID:       1,
		Name:     "Test Book",
		SellerID: 1,
		Costs:    29.99,
	}

	mockRepo.On("GetByID", uint(1)).Return(expectedBook, nil)

	result, err := service.GetBookByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBook.Name, result.Name)
	assert.Equal(t, expectedBook.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetBookByID_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetBookByID(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetBookByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("database error"))

	result, err := service.GetBookByID(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetBooksBySellerID_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	expectedBooks := []model.Book{
		{ID: 1, Name: "Book 1", SellerID: 1, Costs: 29.99},
		{ID: 2, Name: "Book 2", SellerID: 1, Costs: 39.99},
	}

	mockRepo.On("GetBySellerID", uint(1)).Return(expectedBooks, nil)

	result, err := service.GetBooksBySellerID(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedBooks[0].Name, result[0].Name)
	assert.Equal(t, expectedBooks[1].Name, result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &model.Book{
		ID:          1,
		SellerID:    1,
		Name:        "Old Name",
		Description: "Old Description",
		Author:      "Old Author",
		Stock:       5,
		Costs:       19.99,
		Category:    "Old Category",
	}

	newName := "New Name"
	newCosts := 29.99
	req := &model.UpdateBookRequest{
		Name:  &newName,
		Costs: &newCosts,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Book")).Return(nil)

	result, err := service.UpdateBook(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newCosts, result.Costs)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook_BookNotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	req := &model.UpdateBookRequest{}

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.UpdateBook(1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook_Unauthorized(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &model.Book{
		ID:       1,
		SellerID: 1,
		Name:     "Test Book",
	}

	req := &model.UpdateBookRequest{}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)

	result, err := service.UpdateBook(1, req, 2)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "unauthorized: you can only update your own books", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &model.Book{
		ID:       1,
		SellerID: 1,
		Name:     "Test Book",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)
	mockRepo.On("Delete", uint(1), uint(1)).Return(nil)

	err := service.DeleteBook(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook_BookNotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	err := service.DeleteBook(1, 1)

	assert.Error(t, err)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook_Unauthorized(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &model.Book{
		ID:       1,
		SellerID: 1,
		Name:     "Test Book",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)

	err := service.DeleteBook(1, 2)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized: you can only delete your own books", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeductStock_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &model.Book{
		ID:       1,
		SellerID: 1,
		Name:     "Test Book",
		Stock:    10,
		Costs:    29.99,
	}

	updatedBook := &model.Book{
		ID:       1,
		SellerID: 1,
		Name:     "Test Book",
		Stock:    7,
		Costs:    29.99,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil).Once()
	mockRepo.On("DeductStock", uint(1), 3).Return(nil)
	mockRepo.On("GetByID", uint(1)).Return(updatedBook, nil).Once()

	result, err := service.DeductStock(1, 3)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 7, result.Stock)
	mockRepo.AssertExpectations(t)
}

func TestDeductStock_InvalidAmount(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	result, err := service.DeductStock(1, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "deduction amount must be greater than 0", err.Error())
}

func TestDeductStock_BookNotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.DeductStock(1, 3)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeductStock_InsufficientStock(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &model.Book{
		ID:       1,
		SellerID: 1,
		Name:     "Test Book",
		Stock:    2,
		Costs:    29.99,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)

	result, err := service.DeductStock(1, 5)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "insufficient stock", err.Error())
	mockRepo.AssertExpectations(t)
}