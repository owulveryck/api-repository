package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/api-repository/business"
	"github.com/owulveryck/api-repository/handler"

	//_ "github.com/owulveryck/api-repository/repository/gcs"
	_ "github.com/owulveryck/api-repository/repository/fs"
	_ "github.com/owulveryck/api-repository/session/memory"
	"github.com/owulveryck/api-repository/worker"
)

type configuration struct {
	MaxWorkers int   `envconfig:"MAX_WORKERS" default:"10"`
	MaxQueue   int   `envconfig:"MAX_QUEUE" default:"100"`
	MaxLength  int64 `envconfig:"MAX_LENGTH" default:"10240"`
	Port       int   `envconfig:"PORT" default:"8080"`
}

var (

	// jobQueue is A buffered channel that we can send work requests on.
	jobQueue chan<- worker.Job
	config   configuration
)

const (
	productPath   = "products/"
	attributePath = "attributes/"
)

func main() {
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	jobQueue = worker.NewJobQueue(config.MaxQueue)
	dispatcher := worker.NewDispatcher(config.MaxWorkers)
	dispatcher.Run()

	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.Handle("/products", handler.SimplePost{
		Element:  &business.Product{},
		Path:     productPath,
		JobQueue: jobQueue,
	})

	log.Printf("Listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil))
}
