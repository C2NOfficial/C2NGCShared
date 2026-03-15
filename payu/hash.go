package payu

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
)

// Reference: https://docs.payu.in/docs/generate-hash-payu-hosted
func GenerateInitialHash(p *PaymentRequest, salt string) string {
	hashString := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|||||||||||%s",
		p.Key, p.TxnID, p.Amount, p.ProductInfo, p.Firstname, p.Email, salt,
	)
	hash := sha512.Sum512([]byte(hashString))
	return hex.EncodeToString(hash[:])
}

// Reference: https://docs.payu.in/docs/generate-hash-payu-hosted
func GenerateResponseHash(form url.Values, salt string) string {
	hashString := fmt.Sprintf(
		"%s|%s||||||%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		salt,
		form.Get("status"),
		form.Get("udf5"),
		form.Get("udf4"),
		form.Get("udf3"),
		form.Get("udf2"),
		form.Get("udf1"),
		form.Get("email"),
		form.Get("firstname"),
		form.Get("productinfo"),
		form.Get("amount"),
		form.Get("txnid"),
		form.Get("key"),
	)
	hash := sha512.Sum512([]byte(hashString))
	return hex.EncodeToString(hash[:])
}