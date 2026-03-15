package delhivery

import "errors"

var(
	ErrInvalidBillingMode        = errors.New("invalid billing mode")
	ErrInvalidChargeableWeight   = errors.New("invalid chargeable weight")
	ErrInvalidOriginPincode      = errors.New("invalid origin pincode")
	ErrInvalidDestinationPincode = errors.New("invalid destination pincode")
	ErrInvalidShipmentStatus     = errors.New("invalid shipment status")
	ErrInvalidPaymentType        = errors.New("invalid payment type")
	ErrInvalidDimensions         = errors.New("invalid dimensions")
	ErrInvalidPackageType        = errors.New("invalid package type")
)