package constants

const (
	CONTACT_US                  = "contact-us"
	GET_INVOICE_PDF             = "get-invoice-pdf"
	PLACE_ORDER                 = "place-order"
	UPDATE_ORDER                = "update-order"
	TRACK_ORDER                 = "track-order"
	ORDER_CLEANUP_CRON_JOB      = "order-cleanup-cron-job"
	PAYU_PAYMENT_CALLBACK       = "payu-payment-callback"
	PAYU_REFUND_WEBHOOK         = "payu-refund-webhook"
	PAYU_RETRY_PAYMENT          = "payu-retry-payment"
	GET_ESTIMATED_SHIPPING_COST = "get-estimated-shipping-cost"
	CREATE_DELHIVERY_SHIPMENT   = "create-delhivery-shipment"
	DELHIVERY_WEBHOOOK          = "delhivery-webhook"
)

type ServiceEndpoint struct {
	Port string
}

var Endpoints = map[string]ServiceEndpoint{
	CONTACT_US: {
		Port: "3002",
	},
	GET_INVOICE_PDF: {
		Port: "3003",
	},
	UPDATE_ORDER: {
		Port: "3004",
	},
	PLACE_ORDER: {
		Port: "3005",
	},
	TRACK_ORDER: {
		Port: "3006",
	},
	PAYU_PAYMENT_CALLBACK: {
		Port: "4000",
	},
	PAYU_REFUND_WEBHOOK: {
		Port: "4002",
	},
	PAYU_RETRY_PAYMENT: {
		Port: "4003",
	},
	GET_ESTIMATED_SHIPPING_COST: {
		Port: "5002",
	},
	CREATE_DELHIVERY_SHIPMENT: {
		Port: "5003",
	},
	DELHIVERY_WEBHOOOK: {
		Port: "5004",
	},
	ORDER_CLEANUP_CRON_JOB: {
		Port: "6000",
	},
}

// Helpful while debugging
func GetEndpointURL(service string) string {
	e, ok := Endpoints[service]
	if !ok {
		panic("Endpoint not found: " + service)
	}
	return "http://localhost:" + e.Port + "/" + service
}
