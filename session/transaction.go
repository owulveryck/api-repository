package session

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/object"
)

// NewTransaction ...
func NewTransaction(o object.Iterator) (*Transaction, uuid.UUID) {
	u := uuid.New()
	t := &Transaction{
		ID:       u.String(),
		Elements: make([]Element, o.Len()),
	}

	o.Reset()
	i := 0
	for o.Next() {
		t.Elements[i] = Element{
			ID:     o.Element().ID(),
			Status: http.StatusAccepted,
		}
		i++
	}
	return t, u
}
