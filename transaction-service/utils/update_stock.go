package utils

import (
	"fmt"
	"io"
	"main/model"
	"net/http"
)

func UpdateStock(trans model.Transaction, qty int, token string) error {
	url := fmt.Sprintf("http://book-service:8081/books/%d/%d", trans.Book_ID, qty)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return ErrBadReq
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // ✅ Add Bearer token here

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ErrBadReq
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update stock: %s", string(bodyBytes))
	}

	return nil
}
