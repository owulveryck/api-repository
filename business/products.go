package business

import (
	"encoding/json"

	"github.com/owulveryck/api-repository/object"
)

// NewProducts ...
func NewProducts() *Products {
	return &Products{
		it: -1,
	}
}

// Products is an array of Product
type Products struct {
	Elements []*Product
	it       int
}

// Next to fulfil the Iterator interface
func (p *Products) Next() bool {
	if p.it+1 < len(p.Elements) {
		p.it++
		return true
	}
	return false
}

// Element ...
func (p *Products) Element() object.IDer {
	return p.Elements[p.it]
}

// Reset ...
func (p *Products) Reset() {
	p.it = -1
}

// MarshalJSON serialize the product in JSON format compatible with Google OMIT
func (p Products) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Elements) // OMIT
}

// UnmarshalJSON deserialize b into p // OMIT
func (p *Products) UnmarshalJSON(b []byte) error {
	// ...
	var elements []*Product
	err := json.Unmarshal(b, &elements) // OMIT
	(*p).Elements = elements
	return err // OMIT
}

// Len of the elements
func (p *Products) Len() int {
	return len(p.Elements)
}
