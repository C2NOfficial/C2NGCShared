package brevo

import (
	"context"
	"log"
	"os"

	gcp_shared "github.com/C2NOfficial/C2NGCShared/gcp"
	brevo_official "github.com/getbrevo/brevo-go/lib"
)

var (
	API_KEY   string
	Cfg       *brevo_official.Configuration
	Client    *brevo_official.APIClient
)

func Init(ctx context.Context, projectID string) {
	//load environment variables
	if os.Getenv("ENV") == "DEBUG" {
		API_KEY = os.Getenv("BREVO_API_KEY")
	} else {
		if projectID == "" {
			log.Fatal("projectID must be provided in production")
		}
		API_KEY = gcp_shared.LoadSecretsHelper(projectID, "BREVO_API_KEY")
	}
	if API_KEY == "" {
		log.Fatal("BREVO_API_KEY was not found in the .env file")
	}
	//Initialize Brevo 
	Cfg = brevo_official.NewConfiguration()
	Cfg.AddDefaultHeader("api-key", API_KEY)
	Client = brevo_official.NewAPIClient(Cfg)

	log.Print("Brevo was successfully initialized")
}
