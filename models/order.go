package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

)

type OrderItem struct {
	ProductId      string                        `json:"productId" firestore:"productId"`
	Name           string                        `json:"name" firestore:"name"`
	HSNCode        string                        `json:"hsnCode" firestore:"hsnCode"`
	Quantity       int                           `json:"quantity" firestore:"quantity"`
	ImageURL       string                        `json:"imageUrl" firestore:"imageUrl"`
	Price          float64                       `json:"price" firestore:"price"`
	SizeDetails    map[string]ProductSizeDetails `json:"sizeDetails" firestore:"sizeDetails"`
	Customizations []CustomizationSelected       `json:"customizations,omitempty" firestore:"customizations,omitempty"`
}

// Returns the first size of the order item.
// Only one size will be present inside the object at any given time.
func (oi *OrderItem) GetSize() string {
	var size string
	for key := range oi.SizeDetails {
		size = key
		break //No need for this since the length will always be one
	}
	return size
}

func (oi *OrderItem) GetEmailFormattedMap() map[string]any {
	return map[string]any{
		"name":           oi.Name,
		"size":           oi.GetSize(),
		"quantity":       oi.Quantity,
		"price":          oi.Price,
		"image":          oi.ImageURL,
		"customizations": oi.GetCustomizationEmailMap(),
	}
}

func (oi *OrderItem) GetCustomizationEmailMap() []map[string]any {
	var customizations []map[string]any
	for _, custom := range oi.Customizations {
		customizations = append(customizations, custom.GetEmailFormattedMap())
	}
	return customizations
}

func (oi *OrderItem) ToFailedStockReset(error string) *FailedStockReset {
	return &FailedStockReset{
		ProductID: oi.ProductId,
		Size:      oi.GetSize(),
		Quantity:  oi.Quantity,
		Reason:    error,
	}
}

type CustomizationSelected struct {
	Id    string           `json:"id" firestore:"id"`
	Name  string           `json:"name" firestore:"name"`
	Type  CustomizatonType `json:"type" firestore:"type"`
	Value string           `json:"value" firestore:"value"`
}

func (cs *CustomizationSelected) GetEmailFormattedMap() map[string]any {
	return map[string]any{
		"id":    cs.Id,
		"name":  cs.Name,
		"type":  cs.Type,
		"value": cs.Value,
	}
}

type CheckoutUserData struct {
	Email string `json:"email" firestore:"email"`
	Phone string `json:"phone" firestore:"phone"`
	Name  string `json:"name" firestore:"name"`
}

// Separate type for order status to avoid any typos
type OrderStatus string

func (o OrderStatus) String() string {
	return string(o)
}

type Order struct {
	ID               string            `firestore:"-"`
	InvoiceNumber    string            `json:"invoiceNumber" firestore:"invoiceNumber"`
	Items            []*OrderItem      `json:"items" firestore:"items"`
	CheckoutUserData *CheckoutUserData `json:"checkoutUserData" firestore:"checkoutUserData"`
	Address          *Address          `json:"address" firestore:"address"`
	Tax              float64           `json:"tax" firestore:"tax"`
	Subtotal         float64           `json:"subtotal" firestore:"subtotal"`
	Total            float64           `json:"total" firestore:"total"`
	Status           OrderStatus       `json:"status" firestore:"status"`
	ShippingFee      float64           `json:"shippingFee" firestore:"shippingFee"`
	PaymentRetryTime time.Time         `firestore:"paymentRetryTime"`
	SearchTokens     []string          `firestore:"searchTokens"`
	Waybill          string            `json:"waybill" firestore:"waybill"`
	CreatedAt        time.Time         `json:"createdAt" firestore:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt" firestore:"updatedAt"`
}

func (o *Order) SetID(id string) {
	o.ID = id
}

func (o *Order) SetInvoiceNumber(invoiceNumber string) {
	o.InvoiceNumber = invoiceNumber
}

func (o *Order) SetShippingFee(shippingFee float64) {
	o.ShippingFee = shippingFee
}

func (o *Order) SetCreatedAt(time time.Time) {
	o.CreatedAt = time
}

func (o *Order) SetUpdatedAt(time time.Time) {
	o.UpdatedAt = time
}

func (o *Order) SetSearchTokens() {
	orderID := strings.ToLower(o.ID)
	customerName := strings.ToLower(o.CheckoutUserData.Name)
	customerEmail := strings.ToLower(o.CheckoutUserData.Email)
	customerPhone := strings.ToLower(o.CheckoutUserData.Phone)

	o.SearchTokens = append(o.SearchTokens,
		orderID,
		customerName,
		customerEmail,
		customerPhone,
	)

	// partial name search
	o.SearchTokens = append(o.SearchTokens, strings.Fields(customerName)...)

	// items
	for _, item := range o.Items {
		o.SearchTokens = append(o.SearchTokens,
			strings.ToLower(item.ProductId),
			strings.ToLower(item.Name),
			strings.ToLower(item.GetSize()),
		)
	}

	// address
	if o.Address != nil {
		o.SearchTokens = append(o.SearchTokens,
			strings.ToLower(o.Address.Line),
			strings.ToLower(o.Address.City),
			strings.ToLower(o.Address.Zip),
		)
	}
}

func (o *Order) SetPaymentRetryTime(time time.Time) {
	o.PaymentRetryTime = time
}

func (o *Order) SetWayBill(waybill string) {
	o.Waybill = waybill
}

func (o *Order) IsRetryTokenExpired() bool {
	return time.Now().UTC().After(o.PaymentRetryTime)
}

func (o *Order) GetCalculatedSubtotal() float64 {
	subtotal := 0.0
	for _, item := range o.Items {
		subtotal += item.Price * float64(item.Quantity)
	}
	return subtotal
}

func (o *Order) GetCalculatedTax() float64 {
	return o.Subtotal * 0.05
}

func (o *Order) GetTotalItemWeight() float64 {
	totalWeight := 0.0
	// Time complexity is O(n) not O(n^2) since item.SizeDetails length is always gonna be 1
	for _, item := range o.Items {
		for _, sizeDetail := range item.SizeDetails {
			totalWeight += sizeDetail.Weight * float64(item.Quantity)
		}
	}
	return totalWeight
}

func (o *Order) GetLocalCreatedAt() time.Time {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Now()
	}
	return o.CreatedAt.In(loc)
}

func (o *Order) GetLocalUpdatedAt() time.Time {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Now()
	}
	return o.UpdatedAt.In(loc)
}

// Returns the formatted created at time in local time using Timezone
//
// Format: "January 02 2006"
func (o *Order) GetFormattedCreatedAtDate() string {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return o.CreatedAt.Local().Format("January 02 2006")
	}
	return o.CreatedAt.In(loc).Format("January 02 2006")
}

// Returns the formatted updated at time in local time using Timezone
//
// Format: "January-02-2006 at 3:04 PM"
func (o *Order) GetFormattedUpdatedAt() string {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return o.UpdatedAt.Local().Format("January-02-2006 at 3:04 PM")
	}
	return o.UpdatedAt.In(loc).Format("January-02-2006 at 3:04 PM")
}

func (o *Order) GetFormattedSubtotal() string {
	return fmt.Sprintf("%.2f", o.Subtotal)
}

func (o *Order) GetFormattedTax() string {
	return fmt.Sprintf("%.2f", o.Tax)
}

func (o *Order) GetFormattedTaxHalf() string {
	return fmt.Sprintf("%.2f", o.Tax/2)
}

func (o *Order) GetFormattedShippingFee() string {
	return fmt.Sprintf("%.2f", o.ShippingFee)
}

func (o *Order) GetFormattedTotal() string {
	return fmt.Sprintf("%.2f", o.Total)
}

func (o *Order) GetItemsEmailFormattedMapList() []map[string]any {
	items := make([]map[string]any, 0)
	for _, item := range o.Items {
		items = append(items, item.GetEmailFormattedMap())
	}
	return items
}

func (o *Order) GetFailedStockResetMap() []map[string]any {
	failedStockResetList := make([]map[string]any, 0)
	for _, item := range o.Items {
		failedStockResetList = append(failedStockResetList, item.GetEmailFormattedMap())
	}
	return failedStockResetList
}

func (o *Order) GetProductIDList() []string {
	productIDList := make([]string, 0)
	for _, item := range o.Items {
		productIDList = append(productIDList, item.ProductId)
	}
	return productIDList
}

func (o *Order) GetFormattedItemCustomizations() []map[string]any {
	customizations := make([]map[string]any, 0)
	for _, item := range o.Items {
		for _, customization := range item.Customizations {
			customizations = append(customizations, customization.GetEmailFormattedMap())
		}
	}
	return customizations
}

func (o *Order) GetProductDesc() string {
	names := []string{}
	for _, item := range o.Items {
		names = append(names, fmt.Sprintf("%s x%d", item.Name, item.Quantity))
	}
	return strings.Join(names, ", ")
}

func (o *Order) GetTotalQuantity() int {
	totalQuantity := 0
	for _, item := range o.Items {
		totalQuantity += item.Quantity
	}
	return totalQuantity
}

func (o *Order) GetHSNCodes() string {
	hsnCodeArray := make([]string, 0)
	for _, item := range o.Items {
		hsnCodeArray = append(hsnCodeArray, item.HSNCode)
	}
	return strings.Join(hsnCodeArray, ",")
}

func (o *Order) ConvertToPDFItems() [][]string {
	//PDF Items are in order
	// #, Description, Hsn Code, Quantity, Unit Price, Total
	pdfItems := make([][]string, 0)
	for i, item := range o.Items {
		customized := ""
		if len(item.Customizations) > 0 {
			customized = " • Customized"
		}
		pdfItem := []string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%s - Size: %s%s", item.Name, item.GetSize(), customized),
			item.HSNCode,
			strconv.Itoa(item.Quantity),
			fmt.Sprintf("₹%.2f", item.Price),
			fmt.Sprintf("₹%.2f", item.Price*float64(item.Quantity)),
		}
		pdfItems = append(pdfItems, pdfItem)
	}
	return pdfItems
}