package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type sessionHandler struct {
	client *datastore.Client
}

func newSessionHandler() (*sessionHandler, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}
	return &sessionHandler{
		client: client,
	}, nil
}

func (s *sessionHandler) HTTPGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		http.Error(w, "Missing argument", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var t transaction
	_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey("Status", id, nil)
		err := tx.Get(key, &t)
		if err != nil {
			return err
		}
		keys := make([]*datastore.Key, len(t.ElementsID))
		for i, t := range t.ElementsID {
			keys[i] = datastore.NameKey("Element", t, key)
		}
		ts := make([]transactionElement, len(t.ElementsID))
		err = tx.GetMulti(keys, ts)
		if err != nil {
			return err
		}
		t.Elements = ts
		return nil
	}, datastore.ReadOnly)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			http.Error(w, "Transaction not found ", http.StatusNotFound)
			return
		}
		http.Error(w, "Cannot fetch transaction "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusMultiStatus)
	enc := json.NewEncoder(w)
	err = enc.Encode(t)
	if err != nil {
		http.Error(w, "Cannot output transaction "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *sessionHandler) Create(ctx context.Context, id uuid.UUID, t *transaction) error {
	t.ElementsID = make([]string, len(t.Elements))
	for i, e := range t.Elements {
		t.ElementsID[i] = e.ID
	}
	key := datastore.NameKey("Status", id.String(), nil)
	_, err := s.client.Put(ctx, key, t)
	if err != nil {
		return err
	}
	keys := make([]*datastore.Key, len(t.ElementsID))
	for i, t := range t.ElementsID {
		keys[i] = datastore.NameKey("Element", t, key)
	}
	_, err = s.client.PutMulti(ctx, keys, t.Elements)
	log.Println(err)
	return err
}
func (s *sessionHandler) Upsert(ctx context.Context, id uuid.UUID, element transactionElement) error {
	_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey("Status", id.String(), nil)
		k := datastore.NameKey("Element", element.ID, key)
		_, err := tx.Put(k, &element)
		return err
	})
	return err
}
