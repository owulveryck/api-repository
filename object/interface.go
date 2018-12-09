package object

import (
	"encoding/json"
)

// IDer is any object that can give its id and serialize itself in JSON
// START_IDER OMIT
type IDer interface {
	ID() string
	json.Marshaler
	json.Unmarshaler
}

// END_IDER OMIT

// Iterator returns a IDer iterator
type Iterator interface {
	// Next advances the iterator and returns whether
	// the next call to the item method will return a
	// non-nil item.
	//
	// Next should be called prior to any call to the
	// iterator's item retrieval method after the
	// iterator has been obtained or reset.
	//
	// The order of iteration is implementation
	// dependent.
	Next() bool
	// Returns the current Element
	Element() IDer
	// Reset returns the iterator to its start position.
	Reset()
	json.Marshaler
	json.Unmarshaler
}
