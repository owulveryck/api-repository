package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/api-repository/business"
	"github.com/owulveryck/api-repository/handler"
	"github.com/owulveryck/api-repository/session"

	_ "github.com/owulveryck/api-repository/dao/gcs"
	_ "github.com/owulveryck/api-repository/session/gds"

	//_ "github.com/owulveryck/api-repository/dao/fs"
	//_ "github.com/owulveryck/api-repository/session/memory"
	"github.com/owulveryck/api-repository/worker"
)

type configuration struct {
	MaxWorkers int   `envconfig:"MAX_WORKERS" default:"5"`
	MaxQueue   int   `envconfig:"MAX_QUEUE" default:"10"`
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

	http.Handle("/product", handler.SimplePost{
		Element:  &business.Product{},
		Path:     productPath,
		JobQueue: jobQueue,
	})
	http.Handle("/products", handler.BundlePost{
		Elements:  business.NewProducts(),
		Path:      productPath,
		JobQueue:  jobQueue,
		MaxLength: config.MaxLength,
	})
	http.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/jobs/")
		u, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "Not an UUID: "+err.Error(), http.StatusBadRequest)
			return
		}
		t, err := session.Get(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		enc := json.NewEncoder(w)
		enc.Encode(*t)
	})

	log.Printf("Listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil))
}
