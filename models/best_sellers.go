package models

type BestSellerItem struct {
	ID    string // `productID_size`
	Count int64  `firestore:"count"`
}
