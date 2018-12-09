package fs

// START_IMPORT OMIT
import (
	// ...
	"context"       // OMIT
	"encoding/json" // OMIT
	"os"            // OMIT
	"path/filepath" // OMIT

	"github.com/kelseyhightower/envconfig"            // OMIT
	"github.com/owulveryck/api-repository/object"     // OMIT
	"github.com/owulveryck/api-repository/repository" // OMIT
)

// END_IMPORT OMIT

type configuration struct {
	Path string `envconfig:"PATH" required:"true" default:"/tmp"`
}

var config configuration

func init() {
	err := envconfig.Process("FS", &config)
	if err != nil {
		panic(err)
	}

	repository.Register(&fsStorage{
		config.Path,
	})
}

// fsStorage implements the Saver interface;
// it encodes and store an object on the filesystem
// START_OBJECT OMIT
type fsStorage struct {
	Path string
}

func (s *fsStorage) Save(ctx context.Context, object object.IDer, path string) error {
	//time.Sleep(200 * time.Millisecond) OMIT
	fpath := filepath.Join(s.Path, path, object.ID())
	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(object)
	return err
}

// END_OBJECT OMIT
