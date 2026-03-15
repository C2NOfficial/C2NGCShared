package invoice

import (
	"fmt"

	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/delhivery"
	"github.com/C2NOfficial/C2NGCShared/models"
)

type PDFData struct {
	InvNo string
	InvDate string
	OrderNo string
	CompanyName string
	CompanyTradeName string
	CompanyWebsite string
	CompanyAddressLine string
	CompanyCity string
	CompanyState string
	CompanyPincode string
	CompanyCountry string
	CompanyGSTIN string
	CompanyContactEmail string
	CompanyInstagram string
	CustomerName string
	CustomerPhone string
	CustomerEmail string
	CustomerAddressLine string
	CustomerCity string
	CustomerState string
	CustomerPincode string
	CustomerCountry string
	Items [][]string
	Subtotal string
	Shipping string
	CGST string
	SGST string
	IGST string
	Total string
}

func NewPDFData(o* models.Order, dfc* delhivery.FirebaseConfig) *PDFData {
	invNumber, IGST := "", ""
	if o.InvoiceNumber != "" {
		invNumber = o.InvoiceNumber
	}
	// If the order is sold within the place of supply state, 
	// then IGST needs to be shown on the order summary. 
	if o.Address.State == dfc.SellerState {
		IGST = fmt.Sprintf("₹%s", o.GetFormattedTax())
	}
	return &PDFData{
		InvNo: invNumber,
		InvDate: o.GetFormattedCreatedAtDate(),
		OrderNo: o.ID,
		CompanyName: constants.CompanyName,
		CompanyTradeName: constants.CompanyTradeName,
		CompanyWebsite: constants.CompanyWebsite,
		CompanyAddressLine: dfc.SellerAdd,
		CompanyCity: dfc.SellerCity,
		CompanyState: dfc.SellerState,
		CompanyPincode: dfc.SellerPin,
		CompanyCountry: o.Address.Country,
		CompanyGSTIN: constants.CompanyGSTIN,
		CompanyContactEmail: constants.CompanyContactEmail,
		CompanyInstagram: constants.CompanyInstagram,
		CustomerName: o.CheckoutUserData.Name,
		CustomerPhone: o.CheckoutUserData.Phone,
		CustomerEmail: o.CheckoutUserData.Email,
		CustomerAddressLine: o.Address.Line,
		CustomerCity: o.Address.City,
		CustomerState: o.Address.State,
		CustomerPincode: o.Address.Zip,
		CustomerCountry: o.Address.Country,
		Items: o.ConvertToPDFItems(),
		Subtotal: fmt.Sprintf("₹%s", o.GetFormattedSubtotal()),
		CGST: fmt.Sprintf("₹%s", o.GetFormattedTaxHalf()),
		SGST: fmt.Sprintf("₹%s", o.GetFormattedTaxHalf()),
		IGST: IGST,
		Shipping: fmt.Sprintf("₹%s", o.GetFormattedShippingFee()),
		Total: fmt.Sprintf("₹%s", o.GetFormattedTotal()),
	}
}