package gcd

// START_IMPORT OMIT
import (
	"context" // OMIT
	// OMIT
	"errors"
	"log" // OMIT

	// OMIT
	"cloud.google.com/go/datastore" // HL
	"github.com/google/uuid"        // OMIT

	// OMIT
	"github.com/kelseyhightower/envconfig"         // OMIT
	"github.com/owulveryck/api-repository/session" // OMIT
	// ...
)

// END_IMPORT OMIT

func init() {
	sessionHandler, err := newSessionHandler()
	if err != nil {
		panic(err)
	}
	session.Register(sessionHandler)
}

type configuration struct {
	ProjectID  string `envconfig:"GCOUD_PROJECT" required:"true"`
	BucketName string `envconfig:"BUCKET" required:"true"`
}

var config configuration

// START_DEF OMIT
type sessionHandler struct {
	client *datastore.Client
}

// END_DEF OMIT

func newSessionHandler() (*sessionHandler, error) {
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}
	return &sessionHandler{
		client: client,
	}, nil
}

// START_GET OMIT
func (s *sessionHandler) Get(ctx context.Context, id uuid.UUID) (*session.Transaction, error) {
	// END_GET OMIT

	var t session.Transaction
	_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey("Status", id.String(), nil)
		err := tx.Get(key, &t)
		if err != nil {
			return err
		}
		keys := make([]*datastore.Key, len(t.ElementsID))
		for i, t := range t.ElementsID {
			keys[i] = datastore.NameKey("Element", t, key)
		}
		ts := make([]session.Element, len(t.ElementsID))
		err = tx.GetMulti(keys, ts)
		if err != nil {
			return err
		}
		t.Elements = ts
		return nil
	}, datastore.ReadOnly)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, errors.New("Transaction not found")
		}
		return nil, err
	}
	return &t, nil
}

// START_CREATE OMIT
func (s *sessionHandler) Create(ctx context.Context, id uuid.UUID, t *session.Transaction) error {
	// END_CREATE OMIT
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

// START_UPSERT OMIT
func (s *sessionHandler) Upsert(ctx context.Context, id uuid.UUID, element session.Element) error {
	// END_UPSERT OMIT
	_, err := s.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey("Status", id.String(), nil)
		k := datastore.NameKey("Element", element.ID, key)
		_, err := tx.Put(k, &element)
		return err
	})
	return err
}
