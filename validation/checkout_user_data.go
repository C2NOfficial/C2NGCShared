package validation

import (
	"net/mail"
	"regexp"

	"github.com/C2NOfficial/C2NGCShared/models"
)

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
var nameRegex = regexp.MustCompile(`^[a-zA-Z\s]{5,50}$`)

func ValidateCheckoutUserData(cud *models.CheckoutUserData) error {
	_, err := mail.ParseAddress(cud.Email)
	if err != nil {
		return ErrInvalidEmail
	}
	if !phoneRegex.MatchString(cud.Phone) {
		return ErrInvalidPhone
	}
	if !nameRegex.MatchString(cud.Name) {
		return ErrInvalidName
	}
	return nil
}