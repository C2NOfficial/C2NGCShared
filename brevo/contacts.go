package brevo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.brevo.com/v3/contacts"
)

func AddToContactList(email string) error {
	payload := AddContactPayload{
		Email:         email,
		ListIds:       []int{6},
		UpdateEnabled: true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", API_KEY)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("brevo API error: status %s, body: %s", resp.Status, string(respBody))
	}

	return nil
}
