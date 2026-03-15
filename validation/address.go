package validation

import (
	"regexp"
	"strings"

	"github.com/C2NOfficial/C2NGCShared/constants"
	"github.com/C2NOfficial/C2NGCShared/models"
)

var zipRegex = regexp.MustCompile(`^[0-9]{6}$`)

func ValidateAddress(address *models.Address) error {
	if address == nil {
		return ErrAddressRequired
	}

	address.Line = strings.TrimSpace(address.Line)
	address.City = strings.TrimSpace(address.City)
	address.State = strings.TrimSpace(address.State)
	address.Zip = strings.TrimSpace(address.Zip)
	address.Country = strings.TrimSpace(address.Country)

	if len(address.City) < 2 || len(address.City) > 50 {
		return ErrInvalidCity
	}
	if !constants.StateNameMap[address.State] {
		return ErrInvalidState
	}
	if address.Country != "India" {
		return ErrUnsupportedCountry
	}
	if !zipRegex.MatchString(address.Zip) {
		return ErrInvalidZip
	}
	return nil
}
