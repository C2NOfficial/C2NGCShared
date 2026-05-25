package gcp_shared

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/pubsub/v2"
)

var (
	PubSubClient   *pubsub.Client
	initPubSubOnce sync.Once
)

func InitPubSub(ctx context.Context, projectIDWithContext string) {
	var err error
	var projectID string
	initPubSubOnce.Do(func() {
		if os.Getenv("ENV") == "DEBUG" {
			projectID = os.Getenv("GCP_PROJECT_ID")
		} else {
			projectID = projectIDWithContext
		}
		PubSubClient, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create PubSub client: %v", err)
		}
		if projectID == "" {
			log.Fatalf("No project id found for the GCP project: %v", err)
		}
		if PubSubClient == nil {
			log.Fatalf("PubSub client not initialized")
		}
	})
	log.Println("PubSub initialized")
}

// Represents the data payload received by a subscriber to a topic. The `Data` bytes
// can then be accessed to decode it to the struct that the `Data` was sent initially
// as marshalled.
type SubMessage struct {
	Message struct {
		Data        []byte            `json:"data"`
		MessageID   string            `json:"messageId"`
		PublishTime string            `json:"publishTime"`
		Attributes  map[string]string `json:"attributes"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// DecodeSubMessageData decodes the message data as a T struct. Returns a pointer of
// type T struct. Using T for future usecases.
func DecodeSubMessageData[T any](msg *SubMessage) (*T, error) {
	var result T
	if err := json.Unmarshal(msg.Message.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PublishMessage publishes a message to a PubSub topic. Returns an error if
// failed to publish the message. PubSubClient must be initialized before calling
// this function.
func PublishMessage(ctx context.Context, data []byte, topicID string) error {
	if PubSubClient == nil {
		return fmt.Errorf("PubSub client not initialized")
	}

	publisher := PubSubClient.Publisher(topicID)
	result := publisher.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
