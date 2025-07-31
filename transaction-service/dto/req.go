package dto

type CreateTransactionRequest struct {
	BookID int `json:"book_id" validate:"required"`
	Qty    int `json:"qty" validate:"required"`
}

type UpdateTransactionStatusRequest struct {
	Qty int `json:"qty" validate:"required"`
}
