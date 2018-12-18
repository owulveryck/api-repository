// +build ignore

package main

// START_IMPORT OMIT
import (
	//OMIT
	"fmt" //OMIT
	"log"
	"time"

	//OMIT
	//OMIT
	"net/http"          //OMIT
	"net/http/httptest" //OMIT
	"os"                //OMIT
	"sync"

	//OMIT
	"github.com/owulveryck/api-repository/business"    // HL
	_ "github.com/owulveryck/api-repository/dao/dummy" // HL
	"github.com/owulveryck/api-repository/injector"

	// END_IMPORT OMIT
	"github.com/owulveryck/api-repository/doc/my"
	"github.com/owulveryck/api-repository/doc/test1/handler"
)

func main() {
	// START_TEST_SERVER OMIT
	http.Handle("/product", handler.SimplePost{
		Element: &business.Product{}, // HL
		Path:    "/product",
	})
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()
	// END_TEST_SERVER OMIT

	wg := new(sync.WaitGroup)
	replyChan := make(chan injector.Reply)
	// START_SEND OMIT
	concurrency := 1                   // HL
	numberOfElements := 10             // HL
	os.Setenv("DUMMY_DURATION", "0ms") // HL
	wg.Add(1)
	go my.SLO.Evaluate(replyChan, numberOfElements, wg)
	concurrencyChan := make(chan struct{}, concurrency) // Number of concurrent calls
	tmpl := `{ "id":"%v", "title":"my title", "description": "description" }`
	t := time.Now()
	for i := 0; i < numberOfElements; i++ {
		concurrencyChan <- struct{}{}
		wg.Add(1)
		go injector.SendPostRequest(ts.URL+"/product", fmt.Sprintf(tmpl, i), concurrencyChan, wg, replyChan)
	}
	wg.Wait()
	log.Println(time.Since(t))
	// END_SEND OMIT
}
