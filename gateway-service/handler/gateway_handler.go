package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type GatewayHandler struct {
	AuthServiceURL        string
	BookServiceURL        string
	TransactionServiceURL string
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		AuthServiceURL:        os.Getenv("AUTH_SERVICE_URL"),
		BookServiceURL:        os.Getenv("BOOK_SERVICE_URL"),
		TransactionServiceURL: os.Getenv("TRANSACTION_SERVICE_URL"),
	}
}

// Auth

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{fullname=string,email=string,password=string,address=string,role=string} true "User registration data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 500 {object} object{message=string}
// @Router /auth/register [post]
func (h *GatewayHandler) Register(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/register")
}
// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "User login credentials"
// @Success 200 {object} object{message=string,data=object{token=string,user=object}}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 500 {object} object{message=string}
// @Router /auth/login [post]
func (h *GatewayHandler) Login(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/login")
}
// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Security BearerAuth
// @Router /auth/users/{id} [get]
func (h *GatewayHandler) GetUserByID(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/"+c.Param("id"))
}
// UpdateUser godoc
// @Summary Update user information
// @Description Update user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer token"
// @Param request body object{fullname=string,email=string,address=string} true "User update data"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Security BearerAuth
// @Router /auth/users/{id} [put]
func (h *GatewayHandler) UpdateUser(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/"+c.Param("id"))
}

// UpdateBalance godoc
// @Summary Update user balance
// @Description Update user account balance
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param Authorization header string true "Bearer token"
// @Param request body object{balance=number} true "Balance update data"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Security BearerAuth
// @Router /auth/users/{id} [patch]
func (h *GatewayHandler) UpdateBalance(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/"+c.Param("id"))
}

// VerifyUser godoc
// @Summary Verify user email
// @Description Verify user email address with token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{token=string} true "Verification token"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Router /auth/users/verify [post]
func (h *GatewayHandler) VerifyUser(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/verify")
}

// ResendVerificationEmail godoc
// @Summary Resend verification email
// @Description Resend email verification to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "User email"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Router /auth/users/resend-verification-email [post]
func (h *GatewayHandler) ResendVerificationEmail(c echo.Context) error {
	return proxyRequest(c, h.AuthServiceURL+"/users/resend-verification-email")
}

// Books

// GetBooks godoc
// @Summary Get all books
// @Description Get list of all books with optional category filter
// @Tags books
// @Accept json
// @Produce json
// @Param category query string false "Book category filter"
// @Success 200 {object} object{message=string,data=array}
// @Failure 500 {object} object{message=string}
// @Router /books [get]
func (h *GatewayHandler) GetBooks(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books")
}

// GetBookByID godoc
// @Summary Get book by ID
// @Description Get detailed information about a specific book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Router /books/{id} [get]
func (h *GatewayHandler) GetBookByID(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books/"+c.Param("id"))
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book listing (seller only)
// @Tags books
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body object{name=string,description=string,author=string,stock=int,costs=number,category=string} true "Book data"
// @Success 201 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Security BearerAuth
// @Router /books [post]
func (h *GatewayHandler) CreateBook(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books")
}

// UpdateBook godoc
// @Summary Update book information
// @Description Update book details (seller only)
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param Authorization header string true "Bearer token"
// @Param request body object{name=string,description=string,author=string,stock=int,costs=number,category=string} true "Book update data"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 403 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Security BearerAuth
// @Router /books/{id} [put]
func (h *GatewayHandler) UpdateBook(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books/"+c.Param("id"))
}

// DeleteBook godoc
// @Summary Delete book
// @Description Delete a book listing (seller only)
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 403 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Security BearerAuth
// @Router /books/{id} [delete]
func (h *GatewayHandler) DeleteBook(c echo.Context) error {
	return proxyRequest(c, h.BookServiceURL+"/books/"+c.Param("id"))
}

// Transactions

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for purchasing a book
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body object{book_id=int,amount=int,name=string} true "Transaction data"
// @Success 201 {object} object{message=string,data=object{transaction_id=string,payment_url=string}}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 500 {object} object{message=string}
// @Security BearerAuth
// @Router /transactions [post]
func (h *GatewayHandler) CreateTransaction(c echo.Context) error {
	return proxyRequest(c, h.TransactionServiceURL+"/transactions")
}

// GetTransactions godoc
// @Summary Get user transactions
// @Description Get all transactions for the authenticated user
// @Tags transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} object{message=string,data=array}
// @Failure 401 {object} object{message=string}
// @Failure 500 {object} object{message=string}
// @Security BearerAuth
// @Router /transactions [get]
func (h *GatewayHandler) GetTransactions(c echo.Context) error {
	return proxyRequest(c, h.TransactionServiceURL+"/transactions")
}

// UpdateTransactionStatus godoc
// @Summary Update transaction status
// @Description Update transaction status (typically used by payment webhook)
// @Tags transactions
// @Accept json
// @Produce json
// @Param trans_id path string true "Transaction ID"
// @Param Authorization header string true "Bearer token"
// @Param request body object{status=string} true "Transaction status update"
// @Success 200 {object} object{message=string,data=object}
// @Failure 400 {object} object{message=string}
// @Failure 401 {object} object{message=string}
// @Failure 404 {object} object{message=string}
// @Failure 500 {object} object{message=string}
// @Security BearerAuth
// @Router /transactions/{trans_id} [put]
func (h *GatewayHandler) UpdateTransactionStatus(c echo.Context) error {
	return proxyRequest(c, h.TransactionServiceURL+"/transactions/"+c.Param("trans_id"))
}

func proxyRequest(c echo.Context, target string) error {

	targetURL := target + "?" + c.QueryParams().Encode()
	fmt.Println("🔗 Forwarding to:", targetURL)

	bodyBuf := new(bytes.Buffer)
	if _, err := io.Copy(bodyBuf, c.Request().Body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	req, err := http.NewRequest(c.Request().Method, targetURL, bodyBuf)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Forward headers (especially Authorization)
	req.Header = make(http.Header)
	for key, values := range c.Request().Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read response from service")
	}

	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
