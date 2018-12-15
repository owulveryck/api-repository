package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"time"

	"github.com/owulveryck/api-repository/business"
	_ "github.com/owulveryck/api-repository/dao/dummy" // HL
	"github.com/owulveryck/api-repository/doc/test1/handler"
	"github.com/owulveryck/api-repository/injector"
)

var ts *httptest.Server

var concurrency, numElements int

func main() {
	log.Println("Starting")
	// START_MAIN OMIT
	http.Handle("/product", handler.SimplePost{
		Element: &business.Product{}, // HL
		Path:    "/product",
	})
	ts = httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()
	// START_MAIN2 OMIT
	os.Setenv("DUMMY_DURATION", "0ms")
	wg := new(sync.WaitGroup)
	concurrency = 200
	numElements = 10000
	durations := make(chan injector.Reply)
	slo := &injector.SLO{
		Latency: map[float64]time.Duration{
			0.95: 400 * time.Millisecond,
			0.99: 750 * time.Millisecond,
		},
		Allowed5xxErrors: 100 - 99.9,
		Verbose:          true,
	}
	wg.Add(1)
	go slo.Evaluate(durations, numElements, wg)
	// END_MAIN OMIT
	ws := make(chan struct{}, concurrency) // Number of concurrent calls // HL
	tmpl := `{ "id":"%v", "title":"my title", "description": "description" }`
	for i := 0; i < numElements; i++ { // Number of elements // HL
		ws <- struct{}{}
		wg.Add(1)
		go injector.SendPostRequest(ts.URL+"/product", fmt.Sprintf(tmpl, i), ws, wg, durations) // HL
	}
	wg.Wait()
}
