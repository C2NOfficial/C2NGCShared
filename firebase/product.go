package firebase_shared

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
)

func GetProductMapFromProductIDs(ctx context.Context, productIDs []string) (map[string]*models.Product, error) {
	productMap := make(map[string]*models.Product)

	docRefs := make([]*firestore.DocumentRef, len(productIDs))
	for i, id := range productIDs {
		docRefs[i] = FirestoreClient.Collection(constants.COLLECTION_NAME_PRODUCTS).Doc(id)
	}

	snapshots, err := FirestoreClient.GetAll(ctx, docRefs)
	if err != nil {
		return nil, err
	}

	for _, snap := range snapshots {
		if !snap.Exists() {
			return nil, errors.New("some product in your cart was not found. Please try again later.")
		}
		var product models.Product
		if err := snap.DataTo(&product); err != nil {
			return nil, err
		}
		product.SetID(snap.Ref.ID)
		// No need to handle error here. Just set category as nil
		category, _ := FetchCategoryByID(ctx, snap.Data()["categoryId"].(string))
		product.SetCategory(category)
		productMap[snap.Ref.ID] = &product
	}
	return productMap, nil
}

func IncreaseStock(ctx context.Context, productID, size string, quantity int) error {
	doc := FirestoreClient.Collection(constants.COLLECTION_NAME_PRODUCTS).Doc(productID)
	_, err := doc.Update(ctx, []firestore.Update{{
		Path:  fmt.Sprintf("sizeDetails.%s.stock", size), // SizeDetails is just a map with string keys as Sizes.
		Value: firestore.Increment(quantity),
	}})
	return err
}