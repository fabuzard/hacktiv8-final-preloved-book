package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdateBalance(user_id int, amount float64) error {
	url := fmt.Sprintf("http://auth-service:8080/users/%d", user_id)

	data := map[string]interface{}{
		"Amount": amount,
	}

	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return ErrBadReq
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return ErrBadReq
	}
	defer resp.Body.Close()

	return nil
}
