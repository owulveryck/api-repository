package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type dummySession struct {
	db *sync.Map
	//db map[uuid.UUID]map[string]transactionElement
}

func (s *dummySession) HTTPGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		http.Error(w, "Missing argument", http.StatusBadRequest)
		return
	}
	t := transaction{
		ID: id,
	}
	u, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Transaction not found ", http.StatusNotFound)
		return
	}

	tr, ok := s.db.Load(u)
	if !ok {
		http.Error(w, "Transaction not found ", http.StatusNotFound)
		return
	}
	tr.(*sync.Map).Range(func(k, v interface{}) bool {
		t.Elements = append(t.Elements, v.(transactionElement))
		return true
	})
	/*
		for _, v := range tr.(map[string]transactionElement) {
			t.Elements = append(t.Elements, v)
		}
	*/
	w.WriteHeader(http.StatusMultiStatus)
	enc := json.NewEncoder(w)
	err = enc.Encode(t)
	if err != nil {
		http.Error(w, "Cannot output transaction "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *dummySession) Create(ctx context.Context, id uuid.UUID, t *transaction) error {
	tmp := new(sync.Map)
	//tmp := make(map[string]transactionElement, len(t.Elements))
	//s.db[id] = make(map[string]transactionElement, len(t.Elements))
	for _, e := range t.Elements {
		tmp.Store(e.ID, e)
		//tmp[e.ID] = e
	}
	s.db.Store(id, tmp)
	return nil
}
func (s *dummySession) Upsert(ctx context.Context, id uuid.UUID, element transactionElement) error {
	t, ok := s.db.Load(id)
	if !ok {
		return errors.New("Cannot find element")
	}
	t.(*sync.Map).Store(element.ID, element)
	//t.(map[string]transactionElement)[element.ID] = element
	s.db.Store(id, t)
	return nil
}
