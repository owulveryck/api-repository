package dummy

import (
	"context" // OMIT
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/api-repository/dao"
	"github.com/owulveryck/api-repository/object" // OMIT
	// OMIT
)

// START_INIT OMIT
type configuration struct {
	T   time.Duration `envconfig:"DURATION" required:"true" default:"0s"`
	Log bool
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

func consoleLog(format string, v ...interface{}) {
	if config.Log {
		log.Printf(format, v...)
	}
}

// dummyStorage implements the Saver interface;
// START_OBJECT OMIT
type dummyStorage struct {
	wait     time.Duration
	duration time.Duration
}

func (s *dummyStorage) Save(ctx context.Context, object object.IDer, path string) error {
	consoleLog("Start Saving: %v/%v", path, object.ID())
	s.wait, _ = time.ParseDuration(os.Getenv("DUMMY_DURATION")) // OMIT
	s.duration += s.wait
	time.Sleep(s.duration)
	consoleLog("Done saving: %v/%v", path, object.ID())
	s.duration -= s.wait
	return nil
}

// END_OBJECT OMIT