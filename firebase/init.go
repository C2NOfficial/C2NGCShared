package firebase_shared

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	App              *firebase.App     //Firestore App
	AuthClient       *auth.Client      //Auth Client
	FirestoreClient  *firestore.Client //Firestore Client
	initFirebaseOnce sync.Once
)

func InitFirebase(keyPath, projectID string) {
	initFirebaseOnce.Do(func() {
		ctx := context.Background()
		var err error

		var appOptions []option.ClientOption
		if keyPath != "" {
			appOptions = append(appOptions, option.WithCredentialsFile(keyPath))
		}

		App, err = firebase.NewApp(ctx, nil, appOptions...)
		if err != nil {
			log.Fatalf("Error occurred initializing Firebase: %v", err)
		}

		AuthClient, err = App.Auth(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Auth client: %v", err)
		}

		FirestoreClient, err = App.Firestore(ctx)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore client: %v", err)
		}

		log.Println("Firebase initialized successfully")
	})
}
