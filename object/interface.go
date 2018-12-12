package object

import (
	"encoding/json"
)

// START_IDER OMIT

// IDer is any object that can give its id and serialize itself in JSON
type IDer interface {
	ID() string
	json.Marshaler
	json.Unmarshaler
}

// END_IDER OMIT

// Iterator returns a IDer iterator
// START_ITERATOR OMIT
type Iterator interface {
	// Next advances the iterator and returns whether // OMIT
	// the next call to the item method will return a // OMIT
	// non-nil item. // OMIT
	// // OMIT
	// Next should be called prior to any call to the // OMIT
	// iterator's item retrieval method after the // OMIT
	// iterator has been obtained or reset. // OMIT
	// // OMIT
	// The order of iteration is implementation // OMIT
	// dependent. // OMIT
	Next() bool
	// Returns the current Element // OMIT
	Element() IDer
	// Reset returns the iterator to its start position. // OMIT
	Reset()
	Len() int
	json.Marshaler
	json.Unmarshaler
}

// END_ITERATOR OMIT
