package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/owulveryck/api-repository/business"
	_ "github.com/owulveryck/api-repository/dao/dummy" // HL
	"github.com/owulveryck/api-repository/doc/test1/handler"
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
	concurrency = 20
	numElements = 1000
	durations := make(chan time.Duration)
	wg.Add(1)
	go slo(durations, numElements, wg)
	// END_MAIN OMIT
	ws := make(chan struct{}, concurrency) // Number of concurrent calls // HL
	tmpl := `{ "id":"%v", "title":"my title", "description": "description" }`
	for i := 0; i < numElements; i++ { // Number of elements // HL
		ws <- struct{}{}
		wg.Add(1)
		go sendRequest(fmt.Sprintf(tmpl, i), ws, wg, durations) // HL
	}
	wg.Wait()
}

func slo(durations <-chan time.Duration, max int, wg *sync.WaitGroup) {
	defer wg.Done()
	var times []time.Duration
	for t := range durations {
		times = append(times, t)
		if len(times) == max {
			break
		}
	}
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	fmt.Printf("Fastest: %v / Slowest: %v\n", times[0], times[len(times)-1])
	fmt.Printf("95%% of Requests are less than %v\n", times[int(float64(len(times))*0.95)])
	fmt.Printf("99%% of Requests are less than %v\n", times[int(float64(len(times))*0.99)])
}

func sendRequest(object string, w chan struct{}, wg *sync.WaitGroup, c chan<- time.Duration) { // HL
	buf := bytes.NewBufferString(object)
	defer wg.Done()
	start := time.Now()
	resp, err := http.Post(ts.URL+"/product", "application/json", buf)
	if err != nil {
		log.Println(err) // OMIT
	}
	defer resp.Body.Close()
	c <- time.Since(start)
	//c <- time.Since(start)
	//fmt.Println(resp.Status) OMIT
	//io.Copy(os.Stdout, resp.Body) OMIT
	<-w
}
