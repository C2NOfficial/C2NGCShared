package brevo

import (
	"time"

	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
	brevo_official "github.com/getbrevo/brevo-go/lib"
	"google.golang.org/genproto/googleapis/firestore/admin/v1"
)

// helper functions for brevo payloads

func customerTo(o *models.Order) []brevo_official.SendSmtpEmailTo {
	return []brevo_official.SendSmtpEmailTo{{Email: o.CheckoutUserData.Email, Name: o.CheckoutUserData.Name}}
}

func adminTo(mail string) []brevo_official.SendSmtpEmailTo {
	return []brevo_official.SendSmtpEmailTo{{Email: mail, Name: constants.CompanyTradeName}}
}

func adminBCC(mail string) []brevo_official.SendSmtpEmailBcc {
	return []brevo_official.SendSmtpEmailBcc{{Email: mail, Name: constants.CompanyTradeName}}
}

func shippingParams(o *models.Order) map[string]any {
	return map[string]any{
		"shipping_name":         o.CheckoutUserData.Name,
		"shipping_address_line": o.Address.Line,
		"shipping_city":         o.Address.City,
		"shipping_state":        o.Address.State,
		"shipping_zip":          o.Address.Zip,
		"shipping_country":      o.Address.Country,
	}
}

func customerParams(o *models.Order) map[string]any {
	return map[string]any{
		"customer_name":  o.CheckoutUserData.Name,
		"customer_email": o.CheckoutUserData.Email,
		"customer_phone": o.CheckoutUserData.Phone,
	}
}

func year() int {
	return time.Now().UTC().Year()
}

//Used to reduce code when creating parameters for payloads since
//shipping parameters are the same for all payloads
//
//Eg: Params: merge(shippingParams(o), map[string]any{})
func merge(maps ...map[string]any) map[string]any {
	result := map[string]any{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// customer emails

func ToSuccessEmailPayload(o *models.Order, adminMail string) *Payload {
	return &Payload{
		To:         customerTo(o),
		TemplateID: 3,
		BCC:        adminBCC(adminMail),
		Params: merge(shippingParams(o), map[string]any{
			"name":         o.CheckoutUserData.Name,
			"order_number": o.ID,
			"order_date":   o.GetFormattedCreatedAtDate(),
			"order_total":  o.GetFormattedTotal(),
			"items":        o.GetItemsEmailFormattedMapList(),
			"subtotal":     o.GetFormattedSubtotal(),
			"shipping":     o.GetFormattedShippingFee(),
			"tax":          o.GetFormattedTax(),
			"current_year": year(),
		}),
	}
}

func ToOrderShippedEmailPayload(o *models.Order) *Payload {
	return &Payload{
		To:         customerTo(o),
		TemplateID: 7,
		Params: merge(shippingParams(o), map[string]any{
			"name":         o.CheckoutUserData.Name,
			"order_number": o.ID,
			"waybill":      o.Waybill,
			"shipped_date": o.GetFormattedUpdatedAt(),
			"items":        o.GetItemsEmailFormattedMapList(),
			"current_year": year(),
		}),
	}
}

func ToOrderDeliveredEmailPayload(o *models.Order) *Payload {
	return &Payload{
		To:         customerTo(o),
		TemplateID: 9,
		Params: merge(shippingParams(o), map[string]any{
			"name":           o.CheckoutUserData.Name,
			"order_number":   o.ID,
			"order_total":    o.GetFormattedTotal(),
			"delivered_date": o.GetFormattedUpdatedAt(),
			"items":          o.GetItemsEmailFormattedMapList(),
			"current_year":   year(),
		}),
	}
}

func ToDeliveryFailedUserEmailPayload(o *models.Order) *Payload {
	return &Payload{
		To:         customerTo(o),
		TemplateID: 10,
		Params: merge(shippingParams(o), map[string]any{
			"name":           o.CheckoutUserData.Name,
			"order_number":   o.ID,
			"waybill":        o.Waybill,
			"attempted_date": o.GetLocalUpdatedAt().Format("02/01/2006"),
			"items":          o.GetItemsEmailFormattedMapList(),
			"current_year":   year(),
		}),
	}
}

func ToOrderRTOUserEmailPayload(o *models.Order) *Payload {
	return &Payload{
		To:         customerTo(o),
		TemplateID: 12,
		Params: merge(shippingParams(o), map[string]any{
			"name":         o.CheckoutUserData.Name,
			"order_number": o.ID,
			"order_date":   o.GetFormattedCreatedAtDate(),
			"order_total":  o.GetFormattedTotal(),
			"items":        o.GetItemsEmailFormattedMapList(),
			"current_year": year(),
		}),
	}
}

func ToRefundInitiatedEmailPayload(o *models.Order, adminMail string, amount float64) *Payload {
	return &Payload{
		To:         customerTo(o),
		BCC:        adminBCC(adminMail),
		TemplateID: 15,
		Params: map[string]any{
			"name":          o.CheckoutUserData.Name,
			"order_number":  o.ID,
			"order_date":    o.GetFormattedCreatedAtDate(),
			"refund_amount": amount,
			"items":         o.GetItemsEmailFormattedMapList(),
			"current_year":  year(),
		},
	}
}

func ToRefundApprovedEmailPayload(o *models.Order, adminMail string, amount float64) *Payload {
	return &Payload{
		To:         customerTo(o),
		BCC:        adminBCC(adminMail),
		TemplateID: 16,
		Params: map[string]any{
			"name":          o.CheckoutUserData.Name,
			"order_number":  o.ID,
			"order_date":    o.GetFormattedCreatedAtDate(),
			"refund_amount": amount,
			"items":         o.GetItemsEmailFormattedMapList(),
			"current_year":  year(),
		},
	}
}

// ── admin emails ──────────────────────────────────────────────────────────────

func ToAlertAdminPaymentSuccessfulUnmarkedEmailPayload(o *models.Order, mail string, dbError string) *Payload {
	return &Payload{
		To:         adminTo(mail),
		TemplateID: 5,
		Params: merge(shippingParams(o), customerParams(o), map[string]any{
			"txn_id":              o.ID,
			"order_total":         o.GetFormattedTotal(),
			"order_date":          o.GetFormattedCreatedAtDate(),
			"db_error":            dbError,
			"items":               o.GetItemsEmailFormattedMapList(),
			"subtotal":            o.GetFormattedSubtotal(),
			"shipping":            o.GetFormattedShippingFee(),
			"tax":                 o.GetFormattedTax(),
			"current_year":        year(),
		}),
	}
}

func ToAlertAdminResetStockEmailPayload(o *models.Order, adminMail string, fiMap []map[string]any) *Payload {
	return &Payload{
		To:         adminTo(adminMail),
		TemplateID: 6,
		Params: merge(customerParams(o), map[string]any{
			"txn_id":        o.ID,
			"order_total":   o.GetFormattedTotal(),
			"order_status":  o.Status.String(),
			"order_date":    o.GetFormattedCreatedAtDate(),
			"failed_items":  fiMap,
			"items":         o.GetItemsEmailFormattedMapList(),
			"current_year":  year(),
		}),
	}
}

func ToPickupRequestCreationFailedAlertEmailPayload(o *models.Order, adminMail string, err string) *Payload {
	return &Payload{
		To:         adminTo(adminMail),
		TemplateID: 8,
		Params: merge(customerParams(o), map[string]any{
			"order_number":  o.ID,
			"waybill":       o.Waybill,
			"order_total":   o.GetFormattedTotal(),
			"shipped_date":  o.GetFormattedUpdatedAt(),
			"error":         err,
			"items":         o.GetItemsEmailFormattedMapList(),
			"current_year":  year(),
		}),
	}
}

func ToDeliveryFailedAdminEmailPayload(o *models.Order, adminMail string) *Payload {
	return &Payload{
		To:         adminTo(adminMail),
		TemplateID: 11,
		Params: merge(shippingParams(o), customerParams(o), map[string]any{
			"order_number":  o.ID,
			"waybill":       o.Waybill,
			"order_total":   o.GetFormattedTotal(),
			"current_year":  year(),
		}),
	}
}

func ToOrderRTOAdminEmailPayload(o *models.Order, adminMail string) *Payload {
	return &Payload{
		To:         adminTo(adminMail),
		TemplateID: 13,
		Params: merge(shippingParams(o), customerParams(o), map[string]any{
			"order_number":  o.ID,
			"waybill":       o.Waybill,
			"order_date":    o.GetFormattedCreatedAtDate(),
			"order_total":   o.GetFormattedTotal(),
			"current_year":  year(),
		}),
	}
}

func ToOrderReturnedToWarehouseEmailPayload(o *models.Order, adminMail string) *Payload {
	return &Payload{
		To:         adminTo(adminMail),
		TemplateID: 14,
		Params: merge(customerParams(o), map[string]any{
			"order_number":  o.ID,
			"waybill":       o.Waybill,
			"order_date":    o.GetFormattedCreatedAtDate(),
			"returned_date": o.GetFormattedUpdatedAt(),
			"order_total":   o.GetFormattedTotal(),
			"items":         o.GetItemsEmailFormattedMapList(),
			"current_year":  year(),
		}),
	}
}

func ToRefundRejectedEmailPayload(o *models.Order, adminMail string, amount float64) *Payload {
	return &Payload{
		To:         adminTo(adminMail),
		TemplateID: 17,
		Params: merge(customerParams(o), map[string]any{
			"order_number":  o.ID,
			"order_date":    o.GetFormattedCreatedAtDate(),
			"refund_amount": amount,
			"current_year":  year(),
		}),
	}
}