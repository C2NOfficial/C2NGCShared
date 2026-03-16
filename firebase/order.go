package firebase_shared

import (
	"context"
	"fmt"

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
