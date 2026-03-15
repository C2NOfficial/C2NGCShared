package delhivery

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
)

type BillingMode string

func (b BillingMode) String() string {
	return string(b)
}

type ShipmentStatus string

func (s ShipmentStatus) String() string {
	return string(s)
}

type PaymentType string

func (p PaymentType) String() string {
	return string(p)
}

type PackageType string

func (p PackageType) String() string {
	return string(p)
}

// Reference: one.delhivery.com/developer-portal/document/b2c/detail/calculate-shipping-cost
// The url keys are in lowercase (json keys)
type ShippingEstimate struct {
	BillingMode        BillingMode    `json:"billingMode"`        // Express (E) or Surface (S)
	ChargeableWeight   float64        `json:"chargeableWeight"`   // In grams
	OriginPincode      string         `json:"originPincode"`      // Pickup pincode
	DestinationPincode string         `json:"destinationPincode"` // Delivery pincode
	ShipmentStatus     ShipmentStatus `json:"shipmentStatus"`     // Delivered, RTO, DTO
	PaymentType        PaymentType    `json:"paymentType"`        // Prepaid or COD
	Length             float64        `json:"length"`             // Package length
	Breadth            float64        `json:"breadth"`            // Package breadth
	Height             float64        `json:"height"`             // Package height
	PackageType        PackageType    `json:"packageType"`        // Box or Flyer
}

func (p *ShippingEstimate) ToUrlValues() url.Values {
	values := url.Values{}
	values.Set("md", string(p.BillingMode))
	values.Set("cgm", fmt.Sprintf("%.2f", p.ChargeableWeight))
	values.Set("o_pin", p.OriginPincode)
	values.Set("d_pin", p.DestinationPincode)
	values.Set("ss", string(p.ShipmentStatus))
	values.Set("pt", string(p.PaymentType))
	values.Set("l", fmt.Sprintf("%.2f", p.Length))
	values.Set("b", fmt.Sprintf("%.2f", p.Breadth))
	values.Set("h", fmt.Sprintf("%.2f", p.Height))
	values.Set("ikpg_type", string(p.PackageType))
	return values
}

// Returns a string which is used a cache key to avoid
// repeatedly calculating the same estimate.
func (p *ShippingEstimate) GetCacheKey() string {
	raw := fmt.Sprintf(
		"%s|%.2f|%s|%s|%s|%s|%.2f|%.2f|%.2f|%s",
		p.BillingMode,
		p.ChargeableWeight,
		p.OriginPincode,
		p.DestinationPincode,
		p.ShipmentStatus,
		p.PaymentType,
		p.Length,
		p.Breadth,
		p.Height,
		p.PackageType,
	)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

// Used by the cloud function since we only need the weight and dest pincode
type ShippingEstimateRequest struct {
	TotalWeight float64 `json:"total_weight"`
	DestPincode string  `json:"dest_pincode"`
}

// Reference: https://one.delhivery.com/developer-portal/v1/execute
//
// Executed the test api to see the response
type ShippingEstimateResponse []ShippingEstimateItem

// Adding all the fields just in case. Usually won't need most of these. Only total amount and tax data is important
type ShippingEstimateItem struct {
	AdhocData       map[string]any `json:"adhoc_data"`
	ChargeAIR       float64        `json:"charge_AIR"`
	ChargeAWB       float64        `json:"charge_AWB"`
	ChargeCCOD      float64        `json:"charge_CCOD"`
	ChargeCNC       float64        `json:"charge_CNC"`
	ChargeCOD       float64        `json:"charge_COD"`
	ChargeCOVID     float64        `json:"charge_COVID"`
	ChargeCWH       float64        `json:"charge_CWH"`
	ChargeDEMUR     float64        `json:"charge_DEMUR"`
	ChargeDL        float64        `json:"charge_DL"`
	ChargeDOCUMENT  float64        `json:"charge_DOCUMENT"`
	ChargeDPH       float64        `json:"charge_DPH"`
	ChargeDTO       float64        `json:"charge_DTO"`
	ChargeE2E       float64        `json:"charge_E2E"`
	ChargeFOD       float64        `json:"charge_FOD"`
	ChargeFOV       float64        `json:"charge_FOV"`
	ChargeFS        float64        `json:"charge_FS"`
	ChargeFSC       float64        `json:"charge_FSC"`
	ChargeINS       float64        `json:"charge_INS"`
	ChargeLABEL     float64        `json:"charge_LABEL"`
	ChargeLM        float64        `json:"charge_LM"`
	ChargeMPS       float64        `json:"charge_MPS"`
	ChargePEAK      float64        `json:"charge_PEAK"`
	ChargePOD       float64        `json:"charge_POD"`
	ChargeQC        float64        `json:"charge_QC"`
	ChargeREATTEMPT float64        `json:"charge_REATTEMPT"`
	ChargeROV       float64        `json:"charge_ROV"`
	ChargeRTO       float64        `json:"charge_RTO"`
	ChargeWOD       float64        `json:"charge_WOD"`
	ChargePickup    float64        `json:"charge_pickup"`
	ChargedWeight   float64        `json:"charged_weight"`
	Divisor         int            `json:"divisor"`
	GrossAmount     float64        `json:"gross_amount"`
	Status          string         `json:"status"`
	TotalAmount     float64        `json:"total_amount"`
	WtRuleID        *string        `json:"wt_rule_id"`
	WtSOPType       string         `json:"wt_sop_type"`
	ZonalCL         *string        `json:"zonal_cl"`
	Zone            string         `json:"zone"`
	TaxData         TaxData        `json:"tax_data"`
}

type TaxData struct {
	CGST             float64 `json:"CGST"`
	IGST             float64 `json:"IGST"`
	SGST             float64 `json:"SGST"`
	KrishiKalyanCess float64 `json:"krishi_kalyan_cess"`
	ServiceTax       float64 `json:"service_tax"`
	SwacchBharatTax  float64 `json:"swacch_bharat_tax"`
}

// Payload when admin marks a shipment as shipped on their end
type CreateShipment struct {
	OrderID              string `json:"order_id"`
	PickupDate           string `json:"pickup_date"`
	PickupTime           string `json:"pickup_time"`
	PickupLocation       string `json:"pickup_location"`
	ExpectedPackageCount int    `json:"expected_package_count"`
}

func (p *CreateShipment) SetPickupLocation(l string) {
	p.PickupLocation = l
}

// Reference: https://one.delhivery.com/developer-portal/document/b2c/detail/pickup-scheduling
func (p *CreateShipment) ToPickupRequestUrlValues() url.Values {
	values := url.Values{}
	values.Set("pickup_date", p.PickupDate)
	values.Set("pickup_time", p.PickupTime)
	values.Set("pickup_location", p.PickupLocation)
	values.Set("expected_package_count", fmt.Sprintf("%d", p.ExpectedPackageCount))
	return values
}

// Reference: https://one.delhivery.com/developer-portal/document/b2c/detail/order-creation
type Order struct {
	Name           string  `json:"name"`              // customer name - mandatory
	Phone          string  `json:"phone"`             // customer phone - mandatory
	Add            string  `json:"add"`               // customer address - mandatory
	City           string  `json:"city"`              // mandatory
	State          string  `json:"state"`             // mandatory
	Pin            string  `json:"pin"`               // pincode - mandatory
	Country        string  `json:"country,omitempty"` // defaults to India if empty
	Order          string  `json:"order"`
	Waybill        string  `json:"waybill"`              // leave empty, Delhivery auto-generates
	OrderDate      string  `json:"order_date"`           // format: YYYY-MM-DD
	PaymentMode    string  `json:"payment_mode"`         // "Prepaid" or "COD"
	TotalAmount    float64 `json:"total_amount"`         // order total
	CODAmount      float64 `json:"cod_amount,omitempty"` // only if COD
	ProductsDesc   string  `json:"products_desc"`        // product description
	Quantity       int     `json:"quantity"`
	Weight         float64 `json:"weight"` // in kg
	ShipmentWidth  float64 `json:"shipment_width,omitempty"`
	ShipmentHeight float64 `json:"shipment_height,omitempty"`
	ShipmentLength float64 `json:"shipment_length,omitempty"` // undocumented but accepted
	SellerGSTTIN   string  `json:"seller_gst_tin"`            // mandatory
	HSNCode        string  `json:"hsn_code"`                  // product HSN code - mandatory
	SellerTIN      string  `json:"seller_tin,omitempty"`
	SellerName     string  `json:"seller_name"`
	SellerAdd      string  `json:"seller_add"`
	SellerInv      string  `json:"seller_inv"`      // invoice number
	SellerInvDate  string  `json:"seller_inv_date"` // format: YYYY-MM-DD
	ReturnName     string  `json:"return_name,omitempty"`
	ReturnPhone    string  `json:"return_phone,omitempty"`
	ReturnAdd      string  `json:"return_add,omitempty"`
	ReturnCity     string  `json:"return_city,omitempty"`
	ReturnState    string  `json:"return_state,omitempty"`
	ReturnPin      string  `json:"return_pin,omitempty"`
	EWaybill       string  `json:"e_waybill,omitempty"`
}

func ToOrder(o *models.Order, fc *FirebaseConfig) Order {
	return Order{
		Name:          o.CheckoutUserData.Name,
		Phone:         o.CheckoutUserData.Phone,
		Add:           o.Address.Line,
		City:          o.Address.City,
		State:         o.Address.State,
		Pin:           o.Address.Zip,
		Order:         o.ID,
		OrderDate:     o.GetLocalCreatedAt().Format("2006-01-02"),
		PaymentMode:   PaymentPrepaid.String(),
		TotalAmount:   o.Total,
		ProductsDesc:  o.GetProductDesc(),
		Quantity:      o.GetTotalQuantity(),
		Weight:        o.GetTotalItemWeight() / 1000, // in kg
		HSNCode:       o.GetHSNCodes(),
		SellerName:    fc.SellerName,
		SellerAdd:     fc.SellerAdd,
		SellerInv:     o.InvoiceNumber,
		SellerGSTTIN:  constants.CompanyGSTIN,
		SellerInvDate: o.GetLocalCreatedAt().Format("2006-01-02"),
		ReturnName:    fc.SellerName,
		ReturnPhone:   fc.SellerPhone,
		ReturnAdd:     fc.SellerAdd,
		ReturnCity:    fc.SellerCity,
		ReturnState:   fc.SellerState,
		ReturnPin:     fc.SellerPin,
	}
}

type OrderRequest struct {
	PickupLocation PickupLocation `json:"pickup_location"`
	Shipments      []Order        `json:"shipments"`
}

type PickupLocation struct {
	Name    string `json:"name"`
	Add     string `json:"add"`
	City    string `json:"city"`
	State   string `json:"state"`
	Pin     string `json:"pin"`
	Phone   string `json:"phone"`
	Country string `json:"country,omitempty"`
}

// Reference: https://one.delhivery.com/developer-portal/document/b2c/detail/order-creation
type OrderCreationResponse struct {
	Success          bool            `json:"success"`
	Rmk              string          `json:"rmk"`
	UploadWbn        string          `json:"upload_wbn"`
	PackageCount     int             `json:"package_count"`
	PrepaidCount     int             `json:"prepaid_count"`
	CodCount         int             `json:"cod_count"`
	CodAmount        float64         `json:"cod_amount"`
	CashPickups      float64         `json:"cash_pickups"`
	CashPickupsCount float64         `json:"cash_pickups_count"`
	PickupsCount     int             `json:"pickups_count"`
	ReplacementCount int             `json:"replacement_count"`
	Packages         []PackageResult `json:"packages"`
}

// Reference: https://one.delhivery.com/developer-portal/document/b2c/detail/order-creation
type PackageResult struct {
	Waybill     string   `json:"waybill"`
	RefNum      string   `json:"refnum"`
	Client      string   `json:"client"`
	Payment     string   `json:"payment"`
	CodAmount   float64  `json:"cod_amount"`
	Status      string   `json:"status"`
	SortCode    string   `json:"sort_code"`
	Serviceable bool     `json:"serviceable"`
	Remarks     []string `json:"remarks"`
}

type FirebaseConfig struct {
	SellerName     string      `json:"sellerName" firestore:"sellerName"`
	SellerAdd      string      `json:"sellerAdd" firestore:"sellerAdd"`
	SellerCity     string      `json:"sellerCity" firestore:"sellerCity"`
	SellerState    string      `json:"sellerState" firestore:"sellerState"`
	SellerPin      string      `json:"sellerPin" firestore:"sellerPin"`
	SellerPhone    string      `json:"sellerPhone" firestore:"sellerPhone"`
	ShipmentLength float64     `json:"shipmentLength" firestore:"shipmentLength"`
	ShipmentWidth  float64     `json:"shipmentWidth" firestore:"shipmentWidth"`
	ShipmentHeight float64     `json:"shipmentHeight" firestore:"shipmentHeight"`
	ReturnName     string      `json:"returnName" firestore:"returnName"`
	ReturnAdd      string      `json:"returnAdd" firestore:"returnAdd"`
	ReturnCity     string      `json:"returnCity" firestore:"returnCity"`
	ReturnState    string      `json:"returnState" firestore:"returnState"`
	ReturnPin      string      `json:"returnPin" firestore:"returnPin"`
	ReturnPhone    string      `json:"returnPhone" firestore:"returnPhone"`
	OriginPincode  string      `json:"originPincode" firestore:"originPincode"`
	BillingMode    BillingMode `json:"billingMode" firestore:"billingMode"` // "E" or "S"
	PackageType    PackageType `json:"packageType" firestore:"packageType"` // "Box" or "Flyer"
}

func (fc *FirebaseConfig) ToShippingEstimate(cw float64, destPincode string) *ShippingEstimate {
	return &ShippingEstimate{
		BillingMode:        fc.BillingMode,
		ChargeableWeight:   cw,
		OriginPincode:      fc.OriginPincode,
		DestinationPincode: destPincode,
		ShipmentStatus:     ShipmentDelivered,
		PaymentType:        PaymentPrepaid,
		Length:             fc.ShipmentLength,
		Breadth:            fc.ShipmentWidth,
		Height:             fc.ShipmentHeight,
		PackageType:        fc.PackageType,
	}
}

func (d *FirebaseConfig) ToPickupLocation() PickupLocation {
	return PickupLocation{
		Name:  d.SellerName,
		Add:   d.SellerAdd,
		City:  d.SellerCity,
		State: d.SellerState,
		Pin:   d.SellerPin,
		Phone: d.SellerPhone,
	}
}

// DelhiveryWebhookPayload is the payload sent by Delhivery to your webhook endpoint
// Reference: https://delhivery-express-api-doc.readme.io/reference/tracking-via-push-api-webhook-1
type WebhookPayload struct {
	Shipment WebhookShipment `json:"Shipment"`
}

type WebhookShipment struct {
	AWB         string `json:"AWB"`         // waybill number
	ReferenceNo string `json:"ReferenceNo"` // your order ID
	PickUpDate  string `json:"PickUpDate"`
	NSLCode     string `json:"NSLCode"`
	Sortcode    string `json:"Sortcode"`
	Status      Status `json:"Status"`
}

type Status struct {
	Status         string `json:"Status"`         // human readable e.g. "Delivered"
	StatusType     string `json:"StatusType"`     // code e.g. "DL"
	StatusDateTime string `json:"StatusDateTime"` // "2019-01-09T17:10:42.767"
	StatusLocation string `json:"StatusLocation"`
	Instructions   string `json:"Instructions"`
}

type TrackingResponse struct {
	ShipmentData []ShipmentData `json:"ShipmentData"`
}

type ShipmentData struct {
	Shipment ShipmentDetail `json:"Shipment"`
}

type ShipmentDetail struct {
	PickUpDate      string      `json:"PickUpDate"`
	Destination     string      `json:"Destination"`
	DestReceiveDate *string     `json:"DestRecieveDate"`
	DeliveryDate    *string     `json:"DeliveryDate"` // set when delivered
	ReferenceNo     string      `json:"ReferenceNo"`  // order ID
	Status          Status      `json:"Status"`
	Scans           []ScanEntry `json:"Scans"`
}

type ScanEntry struct {
	ScanDetail ScanDetail `json:"ScanDetail"`
}

type ScanDetail struct {
	ScanDateTime    string `json:"ScanDateTime"`
	ScanType        string `json:"ScanType"` // "UD", "PU", "IT" etc
	Scan            string `json:"Scan"`     // human readable
	StatusDateTime  string `json:"StatusDateTime"`
	ScannedLocation string `json:"ScannedLocation"`
	Instructions    string `json:"Instructions"`
	StatusCode      string `json:"StatusCode"`
}
