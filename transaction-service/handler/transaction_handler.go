package handler

import (
	"fmt"
	"main/dto"
	"main/helper"
	"main/model"
	"main/service"
	"main/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/veritrans/go-midtrans"
)

type TransactionHandler struct {
	serv service.TransactionService
}

func NewTransactionHandler(serv service.TransactionService) *TransactionHandler {
	return &TransactionHandler{serv: serv}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	user_id := c.Get("user_id").(int)
	name := c.Get("name").(string)
	email := c.Get("email").(string)

	var req dto.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrBadReq
	}

	// Get token for calling book-service
	token := c.Request().Header.Get("Authorization")

	// Fetch book data
	book, err := utils.GetBookByID(uint(req.BookID), token)
	if err != nil || book.ID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "book not found")
	}
	fmt.Println("Book Response:", book)

	// Validate stock
	if req.Qty <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "quantity must be greater than 0")
	}
	if req.Qty > book.Stock {
		return echo.NewHTTPError(http.StatusBadRequest, "quantity exceeds available stock")
	}

	// Compute amount
	amount := float64(req.Qty) * book.Cost

	// Build transaction model
	t := model.Transaction{
		Book_ID: req.BookID,
		Amount:  amount,
	}

	// Store transaction
	trans, err := h.serv.CreateTransaction(user_id, t)
	if err != nil {
		return err
	}

	// Generate midtrans payment link
	orderId := fmt.Sprintf("%d-%d", trans.Transaction_ID, time.Now().Unix())
	tokenUrl := utils.MidtransPayment(orderId, int(trans.Amount), name, email)

	// Check if midtrans returned a valid token
	if tokenUrl.Token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, tokenUrl.ErrorMessages)
	}

	// Build response
	res := struct {
		TokenUrl       midtrans.SnapResponse `json:"token_url"`
		Transaction_id int                   `json:"transaction_id"`
	}{
		TokenUrl:       tokenUrl,
		Transaction_id: int(trans.Transaction_ID),
	}

	resp := helper.RespHelper("Transaction created successfully", res)
	return c.JSON(http.StatusCreated, resp)
}

func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	user_id := c.Get("user_id").(int)
	transactions, err := h.serv.GetTransaction(user_id)
	if err != nil {
		return err
	}

	resp := helper.RespHelper("Transactions retrieved successfully", transactions)
	return c.JSON(http.StatusOK, resp)
}

func (h *TransactionHandler) UpdateTransactionStatus(c echo.Context) error {
	transaction_id := c.Param("trans_id")

	var qty dto.UpdateTransactionStatusRequest
	if err := c.Bind(&qty); err != nil {
		return utils.ErrBadReq
	}

	token := c.Request().Header.Get("Authorization")

	trans_id, err := strconv.Atoi(transaction_id)
	if err != nil {
		return utils.ErrBadReq
	}
	trans, err := h.serv.UpdateTransactionStatus(trans_id)
	if err != nil {
		return err
	}

	err = utils.UpdateBalance(trans.User_ID, trans.Amount)
	if err != nil {
		return err
	}
	err = utils.EmailTransaction(trans)
	if err != nil {
		return err
	}

	err = utils.UpdateStock(trans, qty.Qty, token)

	resp := helper.RespHelper("Transaction status updated successfully", trans)
	return c.JSON(http.StatusOK, resp)
}
func (h *TransactionHandler) HandleWebhook(c echo.Context) error {
	// Bind request body to get transaction ID
	var req dto.WebhookRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrBadReq
	}

	// Get transaction by ID using service
	transactions, err := h.serv.GetTransactionByID(int(req.TransactionID))
	if err != nil {
		return utils.ErrBadReq
	}

	if transactions.Status == "success" {
		return utils.ErrBadReq
	}

	// get book data using book ID
	book, err := utils.GetBookByID(uint(transactions.Book_ID), c.Request().Header.Get("Authorization"))
	if err != nil || book.ID == 0 {
		return utils.ErrBadReq
	}

	// change status to success
	updatedTransaction, err := h.serv.UpdateTransactionStatus(int(req.TransactionID))
	if err != nil {
		return utils.ErrBadReq
	}

	// Update stock in book-service
	rawToken := c.Request().Header.Get("Authorization") // "Bearer ey..."
	token := strings.TrimPrefix(rawToken, "Bearer ")    // ✅ only the JWT
	if err := utils.UpdateStock(transactions, req.Qty, token); err != nil {
		return utils.ErrBadReq
	}

	// Update seller balance

	if err := utils.UpdateBalance(int(book.SellerID), transactions.Amount); err != nil {
		return utils.ErrBadReq
	}

	if err := utils.EmailTransaction(updatedTransaction); err != nil {
		return utils.ErrBadReq
	}

	// Return the transaction data for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Transaction completed successfully",
		"status":  true,
		"data": map[string]interface{}{
			"transaction": updatedTransaction,
			"book":        book,
			"email_sent":  true,
		},
	})
}
