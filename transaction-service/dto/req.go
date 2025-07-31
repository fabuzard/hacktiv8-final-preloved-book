package dto

type CreateTransactionRequest struct {
	BookID int `json:"book_id" validate:"required"`
	Qty    int `json:"qty" validate:"required"`
}
type WebhookRequest struct {
	TransactionID uint `json:"transaction_id"`
	Qty           int  `json:"qty"`
}
