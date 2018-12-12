package business

import (
	"encoding/json"
	"errors"
	"net/url"
)

// START_PRODUCT OMIT

// A Product that your potential customers would be searching for on Google
type Product struct {
	SKU                 string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Link                *url.URL `json:"link"`
	ImageLink           *url.URL `json:"image_link"`
	AdditionalImageLink *url.URL `json:"additional_image_link,omitempty"`
	MobileLink          *url.URL `json:"mobile_link,omitempty"`
}

// ID of the product OMIT
func (p *Product) ID() string {
	return p.SKU
}

// MarshalJSON serialize the product in JSON format compatible with Google OMIT
func (p Product) MarshalJSON() ([]byte, error) {
	// ...
	type alias Product            // OMIT
	return json.Marshal(alias(p)) // OMIT
}

// UnmarshalJSON deserialize b into p // OMIT
func (p *Product) UnmarshalJSON(b []byte) error {
	// ...
	type alias Product             // OMIT
	var aux alias                  // OMIT
	err := json.Unmarshal(b, &aux) // OMIT
	*p = (Product)(aux)            // OMIT
	if p.SKU == "" {               // OMIT
		return errors.New("id is null or absent") // OMIT
	} // OMIT
	return err // OMIT
}

// END_PRODUCT OMIT
