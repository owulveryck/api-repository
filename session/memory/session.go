package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/session"
)

func init() {
	session.Register(&dummySession{
		db: new(sync.Map),
	})
}

// START_DEF OMIT
type dummySession struct {
	db *sync.Map
}

// END_DEF OMIT

// START_GET OMIT
func (s *dummySession) Get(ctx context.Context, u uuid.UUID) (*session.Transaction, error) {
	// END_GET OMIT
	t := session.Transaction{
		ID: u.String(),
	}
	tr, ok := s.db.Load(u)
	if !ok {
		return nil, errors.New("Not found")
	}
	tr.(*sync.Map).Range(func(k, v interface{}) bool {
		t.Elements = append(t.Elements, v.(session.Element))
		return true
	})
	return &t, nil
}

// START_CREATE OMIT
func (s *dummySession) Create(ctx context.Context, id uuid.UUID, t *session.Transaction) error {
	// END_CREATE OMIT
	tmp := new(sync.Map)
	for _, e := range t.Elements {
		tmp.Store(e.ID, e)
	}
	s.db.Store(id, tmp)
	return nil
}

// START_UPSERT OMIT
func (s *dummySession) Upsert(ctx context.Context, id uuid.UUID, element session.Element) error {
	// END_UPSERT OMIT
	t, ok := s.db.Load(id)
	if !ok {
		return errors.New("Cannot find element")
	}
	t.(*sync.Map).Store(element.ID, element)
	s.db.Store(id, t)
	return nil
}
