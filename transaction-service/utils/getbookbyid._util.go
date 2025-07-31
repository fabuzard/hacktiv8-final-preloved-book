// utils/book_client.go
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"main/dto"

	"net/http"
)

func GetBookByID(bookID uint, token string) (dto.BookResponse, error) {
	var result dto.GetBookByIDResponse

	url := fmt.Sprintf("http://book-service:8081/books/%d", bookID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result.Data, err
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result.Data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result.Data, fmt.Errorf("book-service returned status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	return result.Data, nil
}
