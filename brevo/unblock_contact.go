package brevo

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func unblockContact(email string) {
	encodedEmail := url.PathEscape(email)
	apiURL := fmt.Sprintf("https://api.brevo.com/v3/smtp/blockedContacts/%s", encodedEmail)

	req, err := http.NewRequest(http.MethodDelete, apiURL, nil)
    if err != nil {
        log.Printf("Failed to create unblock request: %v", err)
        return
    }

    req.Header.Set("api-key", API_KEY)
    req.Header.Set("accept", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("Failed to unblock contact: %v", err)
        return
    }
    defer resp.Body.Close()

    // 204 = successfully unblocked
    // 404 = contact wasn't blocked 
    log.Printf("Unblock status for %s: %d", email, resp.StatusCode)
}