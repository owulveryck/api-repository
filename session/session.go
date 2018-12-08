package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var engine Saver

// Register a session engine
func Register(s Saver) {
	engine = s
}

// Element of a transaction
type Element struct {
	ID     string
	Status int
	Err    string
}

// Transaction ...
type Transaction struct {
	ID         string    `json:"transaction"`
	ElementsID []string  `json:"-"`
	Elements   []Element `json:"elements" datastore:"-"`
	LastUpdate time.Time `datastore:"-"`
}

// Saver is a session saver
type Saver interface {
	Get(context.Context, uuid.UUID) (*Transaction, error)
	Create(context.Context, uuid.UUID, *Transaction) error
	Upsert(context.Context, uuid.UUID, Element) error
}

// Create a transaction with the session is
func Create(ctx context.Context, u uuid.UUID, t *Transaction) error {
	return engine.Create(ctx, u, t)
}

// Get ...
func Get(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	return engine.Get(ctx, id)
}

// Upsert ...
func Upsert(ctx context.Context, u uuid.UUID, e Element) error {
	return engine.Upsert(ctx, u, e)
}
