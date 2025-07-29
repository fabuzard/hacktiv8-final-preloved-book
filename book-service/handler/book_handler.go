package handler

import (
	"book-service/model"
	"book-service/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) CreateBook(c echo.Context) error {
	var req model.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	sellerIDStr := c.Get("user_id").(string)
	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid seller ID",
		})
	}

	book, err := h.bookService.CreateBook(&req, uint(sellerID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Book created successfully",
		"data":    book,
	})
}

func (h *BookHandler) GetAllBooks(c echo.Context) error {
	category := c.QueryParam("category")
	
	books, err := h.bookService.GetAllBooks(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Books retrieved successfully",
		"data":    books,
	})
}

func (h *BookHandler) GetBookByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid book ID",
		})
	}

	book, err := h.bookService.GetBookByID(uint(id))
	if err != nil {
		if err.Error() == "book not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Book not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Book retrieved successfully",
		"data":    book,
	})
}

func (h *BookHandler) GetMyBooks(c echo.Context) error {
	sellerIDStr := c.Get("user_id").(string)
	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid seller ID",
		})
	}

	books, err := h.bookService.GetBooksBySellerID(uint(sellerID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "My books retrieved successfully",
		"data":    books,
	})
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid book ID",
		})
	}

	var req model.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	sellerIDStr := c.Get("user_id").(string)
	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid seller ID",
		})
	}

	book, err := h.bookService.UpdateBook(uint(id), &req, uint(sellerID))
	if err != nil {
		if err.Error() == "book not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Book not found",
			})
		}
		if err.Error() == "unauthorized: you can only update your own books" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Forbidden: You can only update your own books",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Book updated successfully",
		"data":    book,
	})
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid book ID",
		})
	}

	sellerIDStr := c.Get("user_id").(string)
	sellerID, err := strconv.ParseUint(sellerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid seller ID",
		})
	}

	err = h.bookService.DeleteBook(uint(id), uint(sellerID))
	if err != nil {
		if err.Error() == "book not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Book not found",
			})
		}
		if err.Error() == "unauthorized: you can only delete your own books" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Forbidden: You can only delete your own books",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Book deleted successfully",
	})
}