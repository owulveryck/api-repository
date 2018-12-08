package gcs

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/api-repository/object"
	"github.com/owulveryck/api-repository/repository"
)

type configuration struct {
	ProjectID string `envconfig:"PROJECT" required:"true"`
	Bucket    string `envconfig:"BUCKER" required:"true"`
}

var config configuration

func init() {
	err := envconfig.Process("GCP", &config)
	if err != nil {
		panic(err)
	}

	s, err := newGCPStorage(context.Background(), config.Bucket)
	if err != nil {
		panic(err)
	}
	repository.Register(s)
}

type gcpStorage struct {
	client *storage.Client
	bkt    *storage.BucketHandle
}

// NewGCPStorage returns a ready to use client.
// It creates the bucket if it does not exists
func newGCPStorage(ctx context.Context, bucketName string) (*gcpStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bkt := client.Bucket(bucketName)
	// Check if the bucket exists
	_, err = bkt.Attrs(ctx)
	if err != nil {
		if err = bkt.Create(ctx, config.ProjectID, nil); err != nil {
			return nil, err
		}
	}
	return &gcpStorage{
		client: client,
		bkt:    client.Bucket(bucketName),
	}, nil
}

// Save to fulfill the interface
func (g *gcpStorage) Save(ctx context.Context, object object.IDer, path string) error {
	obj := g.bkt.Object(path + object.ID())
	w := obj.NewWriter(ctx)
	enc := json.NewEncoder(w)
	err := enc.Encode(object)
	if err != nil {
		return err
	}
	return w.Close()
}
