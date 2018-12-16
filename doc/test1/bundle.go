package main

// START_IMPORT OMIT
import (
	//OMIT
	"fmt" //OMIT

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
	http.Handle("/product", handler.BundlePost{
		Elements: &business.Products{}, // HL
		Path:     "/product",
	})
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()
	// END_TEST_SERVER OMIT

	wg := new(sync.WaitGroup)
	replyChan := make(chan injector.Reply)
	// START_SEND OMIT
	concurrency := 1                   // HL
	numberOfElements := 10             // HL
	numberOfElementsInPayload := 10    // HL
	os.Setenv("DUMMY_DURATION", "0ms") // HL
	wg.Add(1)
	go my.SLO.Evaluate(replyChan, numberOfElements, wg)
	concurrencyChan := make(chan struct{}, concurrency) // Number of concurrent calls
	for i := 0; i < numberOfElements; i++ {
		concurrencyChan <- struct{}{}
		wg.Add(1)
		// Construct the payload ...
		elements := constructPayload(i, numberOfElementsInPayload)
		go injector.SendPostRequest(ts.URL+"/product", elements, concurrencyChan, wg, replyChan)
	}
	wg.Wait()
	// END_SEND OMIT
}

func constructPayload(i, numberOfElementsInPayload int) string {
	tmpl := `{ "id":"%v_%v", "title":"my title", "description": "description" }`
	elements := `[`                                  // OMIT
	separator := `,`                                 // OMIT
	for j := 0; j < numberOfElementsInPayload; j++ { // OMIT
		elements = fmt.Sprintf("%v%v%v", // OMIT
			elements,                // OMIT
			fmt.Sprintf(tmpl, i, j), // OMIT
			separator,               // OMIT
		) // OMIT
		if j == numberOfElementsInPayload-2 { // OMIT
			separator = `]` // OMIT
		} // OMIT
	} // OMIT
	return elements
}
