package repository

import (
	"book-service/model"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *model.Book) error
	GetAll(category string) ([]model.Book, error)
	GetByID(id uint) (*model.Book, error)
	GetBySellerID(sellerID uint) ([]model.Book, error)
	Update(book *model.Book) error
	Delete(id uint, sellerID uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) GetAll(category string) ([]model.Book, error) {
	var books []model.Book
	query := r.db.Model(&model.Book{})
	
	if category != "" {
		query = query.Where("category ILIKE ?", "%"+category+"%")
	}
	
	err := query.Find(&books).Error
	return books, err
}

func (r *bookRepository) GetByID(id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetBySellerID(sellerID uint) ([]model.Book, error) {
	var books []model.Book
	err := r.db.Where("seller_id = ?", sellerID).Find(&books).Error
	return books, err
}

func (r *bookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint, sellerID uint) error {
	return r.db.Where("id = ? AND seller_id = ?", id, sellerID).Delete(&model.Book{}).Error
}