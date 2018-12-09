package dummy

import (
	"context" // OMIT
	"time"

	"github.com/owulveryck/api-repository/object"     // OMIT
	"github.com/owulveryck/api-repository/repository" // OMIT
)

func init() {
	repository.Register(&dummyStorage{
		wait: 200 * time.Millisecond,
	})
}

// dummyStorage implements the Saver interface;
// START_OBJECT OMIT
type dummyStorage struct {
	wait time.Duration
}

func (s *dummyStorage) Save(ctx context.Context, object object.IDer, path string) error {
	for {
		select {
		case <-time.After(s.wait):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// END_OBJECT OMIT
