package firebase_shared

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
)

func FetchOrderFromFirestore(ctx context.Context, orderId string) (*models.Order, error) {
	var order *models.Order

	docSnapshot, err := FirestoreClient.Collection(constants.COLLECTION_NAME_ORDERS).Doc(orderId).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error occurred while fetching order: %v", err)
	}
	err = docSnapshot.DataTo(&order)
	order.SetID(docSnapshot.Ref.ID)
	if err != nil {
		return nil, fmt.Errorf("Error occurred while parsing original order: %v", err)
	}
	return order, nil
}

// Reset the stock for an order.
//
// Returns:
// - []*models.FailedStockReset: List of failed stock resets. This can be converted to
// a map with a getter method on this struct. Used mostly when sending mails with each
// item error if any failed.
func HandleStockReset(ctx context.Context, order *models.Order) []*models.FailedStockReset {
	var stockResetErrors []*models.FailedStockReset
	for _, item := range order.Items {
		if err := IncreaseStock(ctx, item.ProductId, item.GetSize(), item.Quantity); err != nil {
			stockResetErrors = append(stockResetErrors, &models.FailedStockReset{
				ProductID: item.ProductId,
				Size:      item.GetSize(),
				Quantity:  item.Quantity,
				Reason:    err.Error(),
			})
		}
	}
	return stockResetErrors
}

func UpdateBestSeller(ctx context.Context, items []*models.OrderItem) []string {
	if len(items) == 0 {
		return nil
	}
	errors := make([]string, 0)
	FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, item := range items {
			bestSellerItem := item.ToBestSellerItem()
			docRef := FirestoreClient.Collection(constants.COLLECTION_NAME_BEST_SELLERS).Doc(bestSellerItem.ID)
			err := tx.Set(docRef, map[string]interface{}{
				"count": firestore.Increment(bestSellerItem.Count),
			}, firestore.MergeAll)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Error updating best seller %s: %v", bestSellerItem.ID, err))
			}
		}
		return nil
	})
	if len(errors) > 0 {
		return errors
	}
	return nil
}