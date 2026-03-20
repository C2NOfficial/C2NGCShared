package constants

const (
	ORDER_STATUS_PENDING          = "pending"
	ORDER_STATUS_PAID             = "paid"
	ORDER_STATUS_FAILED           = "failed"
	ORDER_STATUS_SHIPPED          = "shipped"
	ORDER_STATUS_DELIVERED        = "delivered"
	ORDER_STATUS_CANCELLED        = "cancelled"
	ORDER_STATUS_REFUND_INITIATED = "refund_initiated"
	ORDER_STATUS_REFUNDED         = "refunded"
	ORDER_STATUS_RTO_DELIVERED    = "rto_delivered"
)

var ORDER_STATUSES_MAP = map[string]bool{
	ORDER_STATUS_PENDING:       true,
	ORDER_STATUS_PAID:          true,
	ORDER_STATUS_FAILED:        true,
	ORDER_STATUS_SHIPPED:       true,
	ORDER_STATUS_DELIVERED:     true,
	ORDER_STATUS_CANCELLED:     true,
	ORDER_STATUS_REFUNDED:      true,
	ORDER_STATUS_RTO_DELIVERED: true,
}
