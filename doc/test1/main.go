package main

// START_IMPORT OMIT
import (
	"bytes"             //OMIT
	"fmt"               //OMIT
	"io"                //OMIT
	"log"               //OMIT
	"net/http"          //OMIT
	"net/http/httptest" //OMIT
	"os"                //OMIT
	"time"              //OMIT

	"github.com/owulveryck/api-repository/business"    // HL
	_ "github.com/owulveryck/api-repository/dao/dummy" // HL

	// END_IMPORT OMIT
	"github.com/owulveryck/api-repository/doc/test1/handler"
)

func main() {
	// START_MAIN2 OMIT
	os.Setenv("DUMMY_DURATION", "0s")
	// START_MAIN OMIT
	http.Handle("/product", handler.SimplePost{
		Element: &business.Product{}, // HL
		Path:    "/product",
	})
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()
	for i := 0; i < 10; i++ { // HL
		buf := bytes.NewBufferString(`{ "id":"1234", "title":"my title", "description": "description" }`)
		start := time.Now()
		resp, err := http.Post(ts.URL+"/product", "application/json", buf) // HL
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Time taken", time.Since(start))
		fmt.Println(resp.Status)
		io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
	} // HL
	// END_MAIN OMIT
}
