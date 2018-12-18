package memory

import (
	"context"

	"github.com/google/uuid"
	"github.com/owulveryck/api-repository/session"
)

func init() {
	session.Register(&dummySession{})
}

type dummySession struct {
}

func (s *dummySession) Get(ctx context.Context, u uuid.UUID) (*session.Transaction, error) {
	return &session.Transaction{}, nil
}

func (s *dummySession) Create(ctx context.Context, id uuid.UUID, t *session.Transaction) error {
	return nil
}
func (s *dummySession) Upsert(ctx context.Context, id uuid.UUID, element session.Element) error {
	return nil
}
