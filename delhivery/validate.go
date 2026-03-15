package delhivery

import "regexp"

var zipRegex = regexp.MustCompile(`^[0-9]{6}$`)

func ValidateShippingCostEstimationPayload(payload *ShippingEstimate) error {
	if !ValidBillingModes[payload.BillingMode] {
		return ErrInvalidBillingMode
	}
	//Note: In case when this function when an order is placed by a customer,
	//the validation for this is not required since each size detail is already
	//validated which the weight.
	if payload.ChargeableWeight <= 0 {
		return ErrInvalidChargeableWeight
	}
	if !zipRegex.MatchString(payload.OriginPincode) {
		return ErrInvalidOriginPincode
	}
	if !zipRegex.MatchString(payload.DestinationPincode) {
		return ErrInvalidDestinationPincode
	}
	if !ValidShipmentStatuses[payload.ShipmentStatus] {
		return ErrInvalidShipmentStatus
	}
	if !ValidPackageTypes[payload.PackageType] {
		return ErrInvalidPackageType
	}
	if payload.Length < 0 || payload.Breadth < 0 || payload.Height < 0 {
		return ErrInvalidDimensions
	}
	if !ValidPackageTypes[payload.PackageType] {
		return ErrInvalidPackageType
	}
	return nil
}
