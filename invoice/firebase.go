package invoice

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/C2NOfficial/C2NGCShared/constants"
	firebase_shared "github.com/C2NOfficial/C2NGCShared/firebase"
)

func GetAndIncrementInvoiceCounter(ctx context.Context) (int64, error) {
    ref := firebase_shared.FirestoreClient.Collection(constants.COLLECTION_NAME_COUNTERS).Doc("invoice")
    var current int64
    err := firebase_shared.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
        doc, err := tx.Get(ref)
        if err != nil {
            return err
        }
        current, _ = doc.Data()["current"].(int64)
        return tx.Update(ref, []firestore.Update{
            {Path: "current", Value: current + 1},
        })
    })
    return current + 1, err
}