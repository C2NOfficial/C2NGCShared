// package gcp is a helper package for shared utility functions required for GCP.
package gcp_shared

import (
	"context"
	"fmt"
	"log"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var (
	secretManagerClient   *secretmanager.Client
	secretManagerErr      error
	initSecretManagerOnce sync.Once
)

func getSecretManagerClient(ctx context.Context) (*secretmanager.Client, error) {
	initSecretManagerOnce.Do(func() {
		secretManagerClient, secretManagerErr = secretmanager.NewClient(ctx)
	})
	return secretManagerClient, secretManagerErr
}

// getSecretFromGCP retrieves the latest version of a secret from Google Cloud Secret Manager.
//
// Parameters:
//   - secretName: The fully-qualified name of the secret version in the format:
//     "projects/{projectID}/secrets/{secretName}/versions/{version}"
//     For example: "projects/my-project/secrets/my-secret/versions/latest"
//
// Returns:
//   - The secret payload as a string if retrieval is successful.
//   - An error if any occurs during client creation or secret access.
func getSecretFromGCP(secretName string) (string, error) {
	ctx := context.Background()

	client, err := getSecretManagerClient(ctx)
	if err != nil {
		return "", err
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}

	secretData := string(result.Payload.Data)
	return secretData, nil
}

func LoadSecretsHelper(projectID string, secretName string) string {
	path := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
	secret, err := getSecretFromGCP(path)
	if err != nil {
		log.Fatalf("Error fetching secret %s: %v", secretName, err)
	}
	return secret
}
