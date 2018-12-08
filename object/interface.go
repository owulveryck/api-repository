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
