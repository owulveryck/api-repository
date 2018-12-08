package main

// AttributeType ...
type AttributeType string

// Attribute ...
type Attribute struct {
	Code string        `json:"code"`
	Name string        `json:"name"`
	Type AttributeType `json:"type,omitempty"`
}

// Model ...
type Model struct {
	Code           string       `json:"code"`
	Name           string       `json:"name"`
	AttributesCode []string     `json:"productAttributeCodes,omitempty"`
	Attributes     []*Attribute `json:"productAttributes,omitempty"`
}

// ID ...
func (m Model) ID() string {
	return m.Code
}

// Products ...
type Products struct {
	Code       string      `json:"code"`
	Model      string      `json:"productModelCode"`
	Attributes []Attribute `json:"productAttributes"`
}

// ReturnCode represent the element of an array in a 207 response
type ReturnCode struct {
	Code    string   `json:"code"`
	Success bool     `json:"success"`
	Details []string `json:"details"`
}
