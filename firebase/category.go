package firebase_shared

import (
	"context"

	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
)

func FetchCategoryByID(ctx context.Context, categoryID string) (*models.Category, error) {
	var category models.Category
	snapshot, err := FirestoreClient.Collection(constants.COLLECTION_NAME_CATEGORIES).Doc(categoryID).Get(ctx)
	if err != nil{
		return nil, err
	}
	err = snapshot.DataTo(&category)
	if err != nil{
		return nil, err
	}
	return &category, nil
}