package validation

import "errors"

var (
	//Common errors for user input validation
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPhone = errors.New("invalid phone")
	ErrInvalidName  = errors.New("invalid name")

	// Address errors
	ErrAddressRequired    = errors.New("address is required")
	ErrInvalidAddressID   = errors.New("invalid address id")
	ErrInvalidAddressLine = errors.New("address line must be between 5 and 70 characters long")
	ErrInvalidCity        = errors.New("city must be at least 2 characters long")
	ErrInvalidState       = errors.New("state must be at least 2 characters long")
	ErrUnsupportedCountry = errors.New("only delivery to India is supported at the moment")
	ErrInvalidZip         = errors.New("zip must be a valid 6-digit Indian pincode")

	//Order errors
	ErrInvalidOrderItems         = errors.New("invalid order items")
	ErrItemsRequired             = errors.New("items are required")
	ErrProductIDRequired         = errors.New("product id is required")
	ErrInvalidOrderStatus        = errors.New("invalid order status")
	ErrInvalidStockQuantity      = errors.New("invalid stock quantity")
	ErrInvalidQuantity           = errors.New("invalid quantity")
	ErrInvalidSize               = errors.New("invalid size")
	ErrSizeRequired              = errors.New("size is required")
	ErrInvalidWeight             = errors.New("invalid weight")
	ErrOutOfStock                = errors.New("out of stock")
	ErrOrderNotFound             = errors.New("order not found")
	ErrIncorrectBackOrderStatus  = errors.New("incorrect back order status")
	ErrInvalidTaxRate            = errors.New("invalid tax rate. Tax rate cannot be 0 or negative")
	ErrInvalidDiscountRate       = errors.New("invalid discount rate. Discount rate cannot be 0 or negative")
	ErrInvalidSubtotal           = errors.New("invalid subtotal")
	ErrInvalidShippingFee        = errors.New("invalid shipping fee")
	ErrInvalidTotal              = errors.New("invalid total")
	ErrInvalidPrice              = errors.New("invalid price")
	ErrInvalidTimeZone           = errors.New("invalid timezone")
	ErrInvalidImageURL           = errors.New("invalid image url")
	ErrInvalidOrderName          = errors.New("invalid order name")
	ErrInvalidProductName        = errors.New("invalid product name")
	ErrInvalidProductSizeDetails = errors.New("invalid product size details")
	ErrInvalidHSNCode            = errors.New("invalid HSN code")
)
