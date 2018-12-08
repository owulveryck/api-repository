package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/owulveryck/api-repository/repository/gcs"
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
	jobQueue chan worker.Job
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
	// Create the storage object
	// Create the dispatcher and launche the worker pools
	dispatcher := worker.NewDispatcher(config.MaxWorkers)
	dispatcher.Run()

	router := httprouter.New()

	// Respond to App Engine and Compute Engine health checks.
	// Indicate the server is healthy.
	router.GET("/_ah/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("ok"))
	})
	router.POST("/products/attributes", HandleAttributeCreate)
	router.POST("/products/models", ProductCreate)
	//  router.POST("/products", Hello)
	router.GET("/products/model_details/:code", HandleProductGet)
	router.GET("/jobs/:id", nil)

	log.Printf("Listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router))
}
