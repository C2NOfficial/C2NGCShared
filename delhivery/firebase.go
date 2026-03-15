package delhivery

import (
	"context"

	"github.com/C2NOfficial/C2NGCShared/constants"
	firebase_shared "github.com/C2NOfficial/C2NGCShared/firebase"
)

func GetFirebaseConfig(ctx context.Context) (*FirebaseConfig, error) {
	docSnapshots, err := firebase_shared.FirestoreClient.Collection(constants.COLLECTION_NAME_DELHIVERY_CONFIG).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var fc FirebaseConfig
	err = docSnapshots[0].DataTo(&fc)
	if err != nil {
		return nil, err
	}
	return &fc, nil
}
