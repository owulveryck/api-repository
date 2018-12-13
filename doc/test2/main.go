package main

import (
	"bytes"
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
)

var ts *httptest.Server

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
	start := time.Now()
	sendRequest := func(object string, w chan struct{}, wg *sync.WaitGroup) { // HL
		buf := bytes.NewBufferString(object)
		defer wg.Done()
		resp, err := http.Post(ts.URL+"/product", "application/json", buf)
		if err != nil {
			log.Println(err) // OMIT
		}
		defer resp.Body.Close()
		//fmt.Println(resp.Status) OMIT
		//io.Copy(os.Stdout, resp.Body) OMIT
		<-w
	}
	wg := new(sync.WaitGroup)
	ws := make(chan struct{}, 10) // Number of concurrent calls // HL
	tmpl := `{ "id":"%v", "title":"my title", "description": "description" }`
	for i := 0; i < 1000; i++ { // Number of elements // HL
		ws <- struct{}{}
		wg.Add(1)
		go sendRequest(fmt.Sprintf(tmpl, i), ws, wg) // HL
	}
	wg.Wait()
	fmt.Println("Time taken", time.Since(start))
	// END_MAIN OMIT
}
