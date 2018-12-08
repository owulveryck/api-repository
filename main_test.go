package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	_ "github.com/owulveryck/api-repository/repository/fs"
)

var testURL string

func TestMain(m *testing.M) {
	config = configuration{
		ProjectID:  "dummy",
		BucketName: "dummy",
		MaxWorkers: 10,
		MaxQueue:   100,
		MaxLength:  15240,
	}

	jobQueue = make(chan Job, config.MaxQueue)
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up
	router := httprouter.New()

	// Respond to App Engine and Compute Engine health checks.
	// Indicate the server is healthy.
	router.GET("/_ah/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("ok"))
	})
	router.POST("/products/attributes", HandleAttributeCreate)
	router.POST("/products/models", ProductCreate)
	router.GET("/products/model_details/:code", HandleProductGet)
	router.GET("/jobs/:id", session.HTTPGet)
	ts := httptest.NewServer(router)
	defer ts.Close()
	testURL = ts.URL

	code := m.Run()
	os.Exit(code)
}
