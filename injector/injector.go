package injector

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// SLO ...
// START_SLO OMIT
type SLO struct {
	Latency          map[float64]time.Duration // Ex: 0.95: 450ms, 0.99: 750ms
	Allowed5xxErrors float64                   // Ex: 0.10 (100% - 99.90%)
	Verbose          bool
}

// END_SLO OMIT

// Evaluate reads from the replyChan til it reach max responses.
// Then it returns the SLO and decrease the waitgroup
// It returns true if the objective is reached.
// It is the responsibility of the call to take care that the max value will be reached
// START_EVALUATE OMIT
// Consuming from replyChan until max is reached and then decrease wg // HL
func (s *SLO) Evaluate(replyChan <-chan Reply, max int, wg *sync.WaitGroup) {
	// END_EVALUATE OMIT
	defer wg.Done()
	var times []time.Duration
	returns := make(map[int]int)
	for r := range replyChan {
		times = append(times, r.timeTaken)
		if len(times) == max {
			break
		}
		returns[r.returnCode]++
	}

	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	slo := true
	if s.Verbose {
		fmt.Printf("Fastest: %v / Slowest: %v\n", times[0], times[len(times)-1])
	}
	for k, v := range s.Latency {
		latency := times[int(float64(len(times))*k)]
		if latency > v {
			slo = false
		}
		if s.Verbose {
			fmt.Printf("%v%% of Requests are less than %v (Expected %v)\n", k*100, latency, v)
		}
	}
	total5xx := 0
	for k, v := range returns {
		if k >= http.StatusInternalServerError {
			total5xx += v
		}
	}
	percent5xx := float64(total5xx) / float64(max) * 100
	if s.Verbose {
		fmt.Printf("%2.2f%% of requests are 5xx (expected %2.2f%%)\n", percent5xx, s.Allowed5xxErrors)
	}
	if percent5xx >= s.Allowed5xxErrors {
		slo = false
	}
	if slo {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}

// Reply from call to SendPostRequest
type Reply struct {
	timeTaken  time.Duration
	returnCode int
	err        error
}

// SendPostRequest a Post Request with the content of object as a Payload.
// START_SEND OMIT
// Sending object to url, send the reply to c; then decrease wg and send event to w once done
func SendPostRequest(url, object string, w chan struct{}, wg *sync.WaitGroup, c chan<- Reply) { // HL
	defer wg.Done()
	start := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(object))
	if err != nil { // OMIT
		c <- Reply{ // OMIT
			time.Since(start),              // OMIT
			http.StatusInternalServerError, // OMIT
			err,                            // OMIT
		} // OMIT
		return // OMIT
	} // OMIT
	defer resp.Body.Close()
	c <- Reply{ // HL
		time.Since(start),
		resp.StatusCode,
		nil,
	}
	//fmt.Println(resp.Status) OMIT
	//io.Copy(os.Stdout, resp.Body) OMIT
	<-w // HL
}

// END_SEND OMIT
