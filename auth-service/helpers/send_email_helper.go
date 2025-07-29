package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendVerificationEmail(email, token string) error {
	emailServiceURL := os.Getenv("EMAIL_SERVICE_URL")
	if emailServiceURL == "" {
		return fmt.Errorf("EMAIL_SERVICE_URL is not set")
	}

	payload := map[string]string{"email": email, "token": token}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(emailServiceURL+"/send-verification-email", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	return nil
}
