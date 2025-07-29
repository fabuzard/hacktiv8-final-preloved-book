package handler

import (
	"main/helper"
	"main/model"
	"main/service"
	"main/utils"
	"net/http"
	"strconv"

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

	var t model.Transaction
	if err := c.Bind(&t); err != nil {
		return utils.ErrBadReq
	}
	trans, err := h.serv.CreateTransaction(user_id, t)

	strId := strconv.Itoa(int(t.Transaction_ID))

	tokenUrl := utils.MidtransPayment(strId, int(trans.Amount), name, email)

	if err != nil {
		return err
	}

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
	user_id := c.Get("id").(int)
	transactions, err := h.serv.GetTransaction(user_id)
	if err != nil {
		return err
	}

	resp := helper.RespHelper("Transactions retrieved successfully", transactions)
	return c.JSON(http.StatusOK, resp)
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
	return c.JSON(http.StatusOK, resp)
}
