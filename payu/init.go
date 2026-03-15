package payu

import (
	"log"
	"os"
	"sync"

	gcp_shared "github.com/C2NOfficial/C2NGCShared/gcp"
)

var (
	Key                string
	Salt               string
	Furl               string
	Surl               string
	PaymentCallbackURL string
	initOnce           sync.Once
)

func Init(projectID string) {
	initOnce.Do(func() {
		if os.Getenv("ENV") == "DEBUG" {
			Key = os.Getenv("PAYU_KEY")
			Salt = os.Getenv("PAYU_SALT")
			Furl = os.Getenv("PAYU_FURL")
			Surl = os.Getenv("PAYU_SURL")
			PaymentCallbackURL = os.Getenv("PAYU_PAYMENT_CALLBACK_URL")
		} else {
			if projectID == "" {
				log.Fatal("projectID must be provided in production")
			}
			Key = gcp_shared.LoadSecretsHelper(projectID, "PAYU_KEY")
			Salt = gcp_shared.LoadSecretsHelper(projectID, "PAYU_SALT")
			Furl = gcp_shared.LoadSecretsHelper(projectID, "PAYU_FURL")
			Surl = gcp_shared.LoadSecretsHelper(projectID, "PAYU_SURL")
			PaymentCallbackURL = gcp_shared.LoadSecretsHelper(projectID, "PAYU_PAYMENT_CALLBACK_URL")
		}
		if Key == "" {
			log.Fatal("PAYU_KEY is not set")
		}
		if Salt == "" {
			log.Fatal("PAYU_SALT is not set")
		}
		if Furl == "" {
			log.Fatal("PAYU_FURL is not set")
		}
		if Surl == "" {
			log.Fatal("PAYU_SURL is not set")
		}
		if PaymentCallbackURL == "" {
			log.Fatal("PAYU_PAYMENT_CALLBACK_URL is not set")
		}
		log.Println("PayU initialized successfully")
	})
}
