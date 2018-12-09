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

type dummySession struct {
	db *sync.Map
}

func (s *dummySession) Get(ctx context.Context, u uuid.UUID) (*session.Transaction, error) {
	/*
		id := ps.ByName("id")
		if id == "" {
			http.Error(w, "Missing argument", http.StatusBadRequest)
			return
		}
		t := session.Transaction{
			ID: id,
		}
		u, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "session.Transaction not found ", http.StatusNotFound)
			return
		}

		tr, ok := s.db.Load(u)
		if !ok {
			http.Error(w, "session.Transaction not found ", http.StatusNotFound)
			return
		}
		tr.(*sync.Map).Range(func(k, v interface{}) bool {
			t.Elements = append(t.Elements, v.(session.TransactionElement))
			return true
		})
		w.WriteHeader(http.StatusMultiStatus)
		enc := json.NewEncoder(w)
		err = enc.Encode(t)
		if err != nil {
			http.Error(w, "Cannot output session.Transaction "+err.Error(), http.StatusInternalServerError)
			return
		}
	*/
	return nil, nil
}

func (s *dummySession) Create(ctx context.Context, id uuid.UUID, t *session.Transaction) error {
	tmp := new(sync.Map)
	for _, e := range t.Elements {
		tmp.Store(e.ID, e)
		//tmp[e.ID] = e
	}
	s.db.Store(id, tmp)
	return nil
}
func (s *dummySession) Upsert(ctx context.Context, id uuid.UUID, element session.Element) error {
	t, ok := s.db.Load(id)
	if !ok {
		return errors.New("Cannot find element")
	}
	t.(*sync.Map).Store(element.ID, element)
	//t.(map[string]session.TransactionElement)[element.ID] = element
	s.db.Store(id, t)
	return nil
}
