package repository

import (
	"context"

	"github.com/owulveryck/api-repository/object"
)

// repository is the object to be used for saving elements
var repository Saver

// Register a storage engine
func Register(s Saver) {
	repository = s
}

// Saver is any object that can saver an IDer
type Saver interface {
	// Save the object into path
	Save(ctx context.Context, object object.IDer, path string) error
}

// Save the object on the given path with the registered engine
func Save(ctx context.Context, object object.IDer, path string) error {
	return repository.Save(ctx, object, path)
}
