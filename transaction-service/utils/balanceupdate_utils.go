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
	_, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))

	if err != nil {
		return ErrBadReq
	}
	return nil
}
