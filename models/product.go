package models

import "time"

type Product struct {
	ID                        string                        `firestore:"-" json:"id"` // Document ID, NOT stored
	Name                      string                        `firestore:"name" json:"name"`
	Slug                      string                        `firestore:"slug" json:"slug"`
	Description               string                        `firestore:"description" json:"description"`
	Category                  *Category                     `firestore:"category" json:"category"`
	Tags                      []string                      `firestore:"tags" json:"tags"`
	Media                     []MediaItem                   `firestore:"media" json:"media"`
	MRP                       float64                       `firestore:"mrp" json:"mrp"`
	CostToMake                float64                       `firestore:"costToMake" json:"costToMake"`
	Quantity                  int                           `firestore:"quantity" json:"quantity"`
	SKU                       string                        `firestore:"sku" json:"sku"`
	AllowOutOfStockPurchase   bool                          `firestore:"allowOutOfStockPurchase" json:"allowOutOfStockPurchase"`
	SizeDetails               map[string]ProductSizeDetails `firestore:"sizeDetails" json:"sizeDetails"`
	Materials                 string                        `firestore:"materials" json:"materials"`
	CareGuide                 string                        `firestore:"careGuide" json:"careGuide"`
	DeliveryPaymentReturnInfo string                        `firestore:"deliveryPaymentReturnInfo" json:"deliveryPaymentReturnInfo"`
	CreatedAt                 time.Time                     `firestore:"createdAt" json:"createdAt"`
	UpdatedAt                 time.Time                     `firestore:"updatedAt" json:"updatedAt"`
}

func (p *Product) SetID(id string) {
	p.ID = id
}

func (p *Product) SetCategory(c *Category) {
	p.Category = c
}

func (p *Product) GetThumbnailImageURL() string {
	return p.Media[0].URL
}

type ProductSizeDetails struct {
	Weight float64 `firestore:"weight" json:"weight"`
	Stock  int     `firestore:"stock" json:"stock"`
}

// When stock reset function is called on each product, this struct is used to represent
// the failed stock reset.
type FailedStockReset struct {
	ProductID string
	Quantity  int
	Size      string
	Reason    string
}

func (fsr *FailedStockReset) toEmailParams() map[string]any {
	return map[string]any{
		"product_id": fsr.ProductID,
		"quantity":   fsr.Quantity,
		"size":       fsr.Size,
		"error":      fsr.Reason,
	}
}

func FSRListToEmailMap(fsrs []*FailedStockReset) []map[string]any {
	emailParams := make([]map[string]any, len(fsrs))
	for i, fsr := range fsrs {
		emailParams[i] = fsr.toEmailParams()
	}
	return emailParams
}
