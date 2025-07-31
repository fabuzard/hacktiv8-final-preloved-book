package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/dto"
	"main/model"
	"net/http"
	"time"
)

func EmailTransaction(trans model.Transaction) error {
	var user dto.GetUserByIDResponse

	urlGetUset := fmt.Sprintf("http://auth-service:8080/users/%d", trans.User_ID)

	req, err := http.NewRequest("GET", urlGetUset, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ErrBadReq
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &user)
	userData := user.User

	url := "http://email-service:8084/send-transaction-success"

	now := time.Now().Format("2006-01-02 15:04:05")

	data := map[string]interface{}{
		"email":          userData.Email,
		"product":        "preloved book",
		"amount":         trans.Amount,
		"transaction_id": trans.Transaction_ID,
		"status":         trans.Status,
		"timestamp":      now,
	}

	jsonData, _ := json.Marshal(data)
	_, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return ErrBadReq
	}

	return nil
}
