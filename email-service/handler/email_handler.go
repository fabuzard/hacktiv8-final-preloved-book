package handler

import (
	"email-service/dto"
	"email-service/utility"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Dummy handler for sending email verification
func SendVerificationEmail(c echo.Context) error {
	var req dto.VerificationEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	// Generate the HTML body with the token
	htmlBody := utility.GenerateVerificationTokenHTML(req.Token)

	// Send the email asynchronously
	go utility.Send(
		[]string{req.Email},
		"Email Verification Code",
		htmlBody,
	)

	// Respond to client
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Verification email sent",
		"email":   req.Email,
	})
}

// Dummy handler for sending transaction success email
func SendTransactionSuccess(c echo.Context) error {
	var req dto.TransactionEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}

	htmlBody := utility.BuildTransactionHTMLBody(
		req.Email,
		req.TransactionID,
		req.Product,
		req.Amount,
		req.Status,
		req.Timestamp,
		req.InvoiceURL,
	)

	err := utility.Send(
		[]string{req.Email},
		"Your Purchase Receipt",
		htmlBody,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Transaction success email sent",
		"email":   req.Email,
	})
}
