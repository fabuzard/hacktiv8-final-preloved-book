package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/model"
	"net/http"
)

func UpdateStock(trans model.Transaction, qty int, token string) error {
	book, err := GetBookByID(uint(trans.Book_ID), token)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://book-service:8081/books/%d", trans.Book_ID)

	data := map[string]interface{}{
		"Stock": book.Stock - qty}

	jsonData, _ := json.Marshal(data)
	_, err = http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return ErrBadReq
	}

	return nil
}
