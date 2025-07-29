package handler

import (
	"main/helper"
	"main/model"
	"main/service"
	"main/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	serv service.TransactionService
}

func NewTransactionHandler(serv service.TransactionService) *TransactionHandler {
	return &TransactionHandler{serv: serv}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	user_id := c.Get("id").(int)
	var t model.Transaction
	if err := c.Bind(&t); err != nil {
		return utils.ErrBadReq
	}
	t, err := h.serv.CreateTransaction(user_id, t)
	if err != nil {
		return err
	}

	resp := helper.RespHelper("Transaction created successfully", t)
	return c.JSON(201, resp)
}

func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	user_id := c.Get("id").(int)
	transactions, err := h.serv.GetTransaction(user_id)
	if err != nil {
		return err
	}

	resp := helper.RespHelper("Transactions retrieved successfully", transactions)
	return c.JSON(200, resp)
}

func (h *TransactionHandler) UpdateTransactionStatus(c echo.Context) error {
	transaction_id := c.Param("trans_id")
	trans_id, err := strconv.Atoi(transaction_id)
	if err != nil {
		return utils.ErrBadReq
	}
	trans, err := h.serv.UpdateTransactionStatus(trans_id)
	if err != nil {
		return err
	}

	resp := helper.RespHelper("Transaction status updated successfully", trans)
	return c.JSON(200, resp)
}
