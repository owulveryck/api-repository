package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/owulveryck/api-repository/business"
	"github.com/owulveryck/api-repository/doc/test1/handler"
	_ "github.com/owulveryck/api-repository/repository/dummy" // HL
)

func main() {
	// START_MAIN2 OMIT
	os.Setenv("DUMMY_DURATION", "0s") // HL
	// START_MAIN OMIT
	http.Handle("/product", handler.SimplePost{
		Element: &business.Product{}, // HL
		Path:    "/product",
	})
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()
	buf := bytes.NewBufferString(`{ "id":"1234", "title":"my title", "description": "description" }`)
	start := time.Now()
	resp, err := http.Post(ts.URL+"/product", "application/json", buf)
	if err != nil {
		log.Println(err)
	}
	log.Println("Time taken", time.Since(start))
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	// END_MAIN OMIT

}
