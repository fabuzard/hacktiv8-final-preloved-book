package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendVerificationEmail(email, token string) error {
	payload := map[string]string{"email": email, "token": token}
	payloadBytes, _ := json.Marshal(payload)

	// change later
	resp, err := http.Post("http://localhost:8084/send-verification-email", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	return nil
}
