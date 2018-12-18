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
	"github.com/owulveryck/api-repository/business"                 // HL
	_ "github.com/owulveryck/api-repository/dao/dummy"              // HL
	_ "github.com/owulveryck/api-repository/doc/test2/session/void" // OMIT
	"github.com/owulveryck/api-repository/injector"
	"github.com/owulveryck/api-repository/worker"

	// END_IMPORT OMIT
	"github.com/owulveryck/api-repository/doc/my"
	"github.com/owulveryck/api-repository/doc/test2/handler"
)

func main() {
	// START_TEST_SERVER OMIT
	jobQueue := worker.NewJobQueue(100)    // HL
	dispatcher := worker.NewDispatcher(50) // HL
	dispatcher.Run()
	http.Handle("/products", handler.BundlePost{
		Elements:  &business.Products{}, // HL
		Path:      "/products",
		JobQueue:  jobQueue,
		MaxLength: 10000,
	})
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()
	wg := new(sync.WaitGroup)
	replyChan := make(chan injector.Reply)
	concurrency := 100                 // HL
	numberOfElements := 10000          // HL
	numberOfElementsInPayload := 100   // HL
	os.Setenv("DUMMY_DURATION", "0ms") // HL
	//...
	// END_TEST_SERVER OMIT
	wg.Add(1)
	go my.SLO.Evaluate(replyChan, numberOfElements, wg)
	concurrencyChan := make(chan struct{}, concurrency) // Number of concurrent calls
	t := time.Now()
	for i := 0; i < numberOfElements; i++ {
		concurrencyChan <- struct{}{}
		wg.Add(1)
		// Construct the payload ...
		elements := constructPayload(i, numberOfElementsInPayload)
		go injector.SendPostRequest(ts.URL+"/products", elements, concurrencyChan, wg, replyChan)
	}
	wg.Wait()
	log.Println(time.Since(t))
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
