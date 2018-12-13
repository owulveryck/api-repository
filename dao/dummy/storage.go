package dummy

import (
	"context" // OMIT
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/api-repository/dao"
	"github.com/owulveryck/api-repository/object" // OMIT
	// OMIT
)

// START_INIT OMIT
type configuration struct {
	T time.Duration `envconfig:"DURATION" required:"true" default:"0s"`
}

var config configuration

func init() {
	err := envconfig.Process("DUMMY", &config)
	// ...
	// END_INIT OMIT
	if err != nil {
		panic(err)
	}
	dao.Register(&dummyStorage{
		wait: config.T,
	})
}

// dummyStorage implements the Saver interface;
// START_OBJECT OMIT
type dummyStorage struct {
	wait time.Duration
}

func (s *dummyStorage) Save(ctx context.Context, object object.IDer, path string) error {
	log.Printf("Start Saving: %v/%v", path, object.ID())
	// for the present tool... OMIT
	envconfig.Process("DUMMY", &config) // OMIT
	s.wait = config.T                   // OMIT
	time.Sleep(s.wait)
	log.Printf("Done saving: %v/%v", path, object.ID())
	return nil
}

// END_OBJECT OMIT
