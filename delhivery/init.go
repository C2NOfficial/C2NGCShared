package delhivery

import (
	"log"
	"os"
	"sync"

	gcp_shared "github.com/C2NOfficial/C2NGCShared/gcp"
)

var (
	API_TOKEN                   string
	WEBHOOK_SECRET              string
	GetEstimatedShippingCostURL string //cloud func url            
	initOnce                    sync.Once
)

func Init(projectID string) {
	initOnce.Do(func() {
		if os.Getenv("ENV") == "DEBUG" {
			API_TOKEN = os.Getenv("DELHIVERY_API_TOKEN")
			WEBHOOK_SECRET = os.Getenv("DELHIVERY_WEBHOOK_SECRET")
			GetEstimatedShippingCostURL = os.Getenv("GET_ESTIMATED_SHIPPING_COST_URL")
		} else {
			if projectID == "" {
				log.Fatal("projectID must be provided in production")
			}
			API_TOKEN = gcp_shared.LoadSecretsHelper(projectID, "DELHIVERY_API_TOKEN")
			WEBHOOK_SECRET = gcp_shared.LoadSecretsHelper(projectID, "DELHIVERY_WEBHOOK_SECRET")
			GetEstimatedShippingCostURL = gcp_shared.LoadSecretsHelper(projectID, "GET_ESTIMATED_SHIPPING_COST_URL")
		}
		if API_TOKEN == "" {
			log.Fatal("DELHIVERY_API_TOKEN is not set")
		}
		if WEBHOOK_SECRET == "" {
			log.Fatal("DELHIVERY_WEBHOOK_SECRET is not set")
		}
		if GetEstimatedShippingCostURL == "" {
			log.Fatal("GET_ESTIMATED_SHIPPING_COST_URL is not set")
		}
		log.Println("Delhivery initialized successfully")
	})
}
