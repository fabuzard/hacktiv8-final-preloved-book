package dto

type VerificationEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}

type TransactionEmailRequest struct {
	Email         string  `json:"email"`
	TransactionID string  `json:"transaction_id"`
	Product       string  `json:"product"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"`
	InvoiceURL    string  `json:"invoice_url"`
}
