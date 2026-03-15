package delhivery

const (
	BillingExpress BillingMode = "E" 
	BillingSurface BillingMode = "S"

	ShipmentDelivered ShipmentStatus = "Delivered"
	ShipmentRTO       ShipmentStatus = "RTO"
	ShipmentDTO       ShipmentStatus = "DTO"

	PaymentPrepaid PaymentType = "Pre-paid"
	PaymentCOD     PaymentType = "COD"

	PackageBox   PackageType = "box"
	PackageFlyer PackageType = "flyer"

	StatusManifested     = "UD"  // order created, not yet picked up
	StatusPickedUp       = "PU"  // courier picked up from warehouse
	StatusInTransit      = "IT"  // moving between hubs
	StatusOutForDelivery = "OD"  // with delivery agent
	StatusDelivered      = "DL"  // delivered to customer
	StatusFailed         = "FA"  // delivery attempt failed
	StatusRTO            = "RTO" // returning to origin
	StatusRTODelivered   = "RTD" // returned to your warehouse

)

// Useful when validating instead of manually checking each
// individually

//
var (
	ValidBillingModes = map[BillingMode]bool{
		BillingExpress: true,
		BillingSurface: true,
	}

	ValidShipmentStatuses = map[ShipmentStatus]bool{
		ShipmentDelivered: true,
		ShipmentRTO:       true,
		ShipmentDTO:       true,
	}

	ValidPaymentTypes = map[PaymentType]bool{
		PaymentPrepaid: true,
		PaymentCOD:     true,
	}

	ValidPackageTypes = map[PackageType]bool{
		PackageBox:   true,
		PackageFlyer: true,
	}
)
