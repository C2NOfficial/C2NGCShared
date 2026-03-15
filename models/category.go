package models

import "time"

type Category struct {
	ID             string                 `firestore:"-" json:"id"` // Firestore document ID
	Name           string                 `firestore:"name" json:"name"`
	HSNCode        string                 `firestore:"hsnCode" json:"hsnCode"`
	Slug           string                 `firestore:"slug" json:"slug"`
	SizeDetails    map[string]SizeDetails `firestore:"sizeDetails" json:"sizeDetails"`
	Customizations []Customization        `firestore:"customizations,omitempty" json:"customizations,omitempty"`
	SizeChart      []map[string][]string  `firestore:"sizeChart" json:"sizeChart"`
	CreatedAt      time.Time              `firestore:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time              `firestore:"updatedAt" json:"updatedAt"`
}

// Used in form of a map[string]SizeDetails where string is the size ('S', 'M', 'L', 'XL'...)
type SizeDetails struct {
	Weight float64 `firestore:"weight" json:"weight"`
}

type Customization struct {
	Id                     string           `firestore:"id" json:"id"`
	Label                  string           `firestore:"label" json:"label"`
	Placeholder            string           `firestore:"placeholder" json:"placeholder"`
	Type                   CustomizatonType `firestore:"type" json:"type"`
	DisabledCustomizations []string         `firestore:"disabledCustomizations,omitempty" json:"disabledCustomizations,omitempty"`
}

type CustomizatonType string

const (
	CustomizationTypeCheckbox CustomizatonType = "checkbox"
	CustomizationTypeText     CustomizatonType = "text"
)
