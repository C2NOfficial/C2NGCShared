package validation

import (
	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
)

func validateOriginalOrderItemsWithPlacedOrderItems(placedOrderItems []*models.OrderItem, originalOrderItems map[string]*models.Product) error {
	if placedOrderItems == nil || len(placedOrderItems) == 0 {
		return ErrItemsRequired
	}
	for _, item := range placedOrderItems {
		original, exists := originalOrderItems[item.ProductId]
		if !exists || original == nil {
			return ErrInvalidOrderItems
		}
		if item.ProductId != original.ID {
			return ErrProductIDRequired
		}
		if item.Name != original.Name {
			return ErrInvalidProductName
		}
		if item.ImageURL != original.GetThumbnailImageURL() {
			return ErrInvalidImageURL
		}
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}
		if item.Price != original.MRP {
			return ErrInvalidPrice
		}
		if item.HSNCode != original.Category.HSNCode {
			return ErrInvalidHSNCode
		}
		//At any given time only one size is allowed in the map since each size is treated as a unique item
		if len(item.SizeDetails) > 1 {
			return ErrInvalidProductSizeDetails
		}
		if err := validateOriginalOrderSizeDetailsWithPlacedOrder(item.SizeDetails, original.SizeDetails); err != nil {
			return err
		}
	}
	return nil
}

func validateOriginalOrderSizeDetailsWithPlacedOrder(psdFromOrder, psdOriginal map[string]models.ProductSizeDetails) error {
	for size, psd := range psdFromOrder {
		//Make sure the key exists in the original map
		original, exists := psdOriginal[size]
		if !exists {
			return ErrInvalidSize
		}
		//Weight should not be changed at all. Very important since delivery fee is calculated based on weight
		if psd.Weight != original.Weight {
			return ErrInvalidWeight
		}
	}
	return nil
}

// Note: This does not includes the shipping fee. Shipping fee is calculated totally separately since
// an additional call needs to be made to the delhivery api to get the original shipping fee.
func validateOrderTotal(o *models.Order, originalProducts map[string]*models.Product) error {
	originalSubTotal := 0.0
	for _, item := range o.Items {
		originalSubTotal += originalProducts[item.ProductId].MRP * float64(item.Quantity)
	}
	if originalSubTotal != o.Subtotal {
		return ErrInvalidSubtotal
	}
	originalTax := o.Subtotal * 0.05
	if originalTax != o.Tax {
		return ErrInvalidTaxRate
	}
	if (originalSubTotal + originalTax) != (o.Subtotal + o.Tax) {
		return ErrInvalidTotal
	}
	return nil
}

func ValidatePlacedOrder(o *models.Order, products map[string]*models.Product) error {
	if o.Status != constants.ORDER_STATUS_PENDING {
		return ErrInvalidOrderStatus
	}
	err := ValidateAddress(o.Address)
	if err != nil {
		return err
	}
	err = ValidateCheckoutUserData(o.CheckoutUserData)
	if err != nil {
		return err
	}
	err = validateOriginalOrderItemsWithPlacedOrderItems(o.Items, products)
	if err != nil {
		return err
	}
	err = validateOrderTotal(o, products)
	if err != nil {
		return err
	}
	return nil
}
