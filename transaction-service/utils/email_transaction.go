package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/model"
	"net/http"
	"time"
)

// Define this locally in transaction-service
type User struct {
	ID         uint      `json:"id"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Password   string    `json:"-"` // ignore for safety
	Address    string    `json:"address"`
	Role       string    `json:"role"`
	Balance    float64   `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
	IsVerified bool      `json:"is_verified"`
}

type GetUserByIDResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

func EmailTransaction(trans model.Transaction) error {
	urlGetUser := fmt.Sprintf("http://auth-service:8080/users/%d", trans.User_ID)

	req, err := http.NewRequest("GET", urlGetUser, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return ErrBadReq
	}
	defer resp.Body.Close()

	var userResp GetUserByIDResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &userResp)
	if err != nil {
		return err
	}
	user := userResp.User

	emailPayload := map[string]interface{}{
		"email":          user.Email,
		"transaction_id": fmt.Sprintf("%d", trans.Transaction_ID),
		"product":        "preloved book",
		"amount":         trans.Amount,
		"status":         trans.Status,
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
		"invoice_url":    "", // blank as requested
	}

	jsonData, err := json.Marshal(emailPayload)
	if err != nil {
		return err
	}

	emailURL := "http://email-service:8084/send-transaction-success"
	req, err = http.NewRequest("POST", emailURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return ErrBadReq
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return ErrBadReq
	}

	return nil
}
