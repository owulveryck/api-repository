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
type SLO struct {
	Latency          map[float64]time.Duration
	Allowed5xxErrors float64
	Verbose          bool
}

// Evaluate reads from the replyChan til it reach max responses.
// Then it returns the SLO and decrease the waitgroup
// It returns true if the objective is reached.
// It is the responsibility of the call to take care that the max value will be reached
func (s *SLO) Evaluate(replyChan <-chan Reply, max int, wg *sync.WaitGroup) {
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
func SendPostRequest(url, object string, w chan struct{}, wg *sync.WaitGroup, c chan<- Reply) { // HL
	buf := bytes.NewBufferString(object)
	defer wg.Done()
	start := time.Now()
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		c <- Reply{
			time.Since(start),
			http.StatusInternalServerError,
			err,
		}
		return
	}
	defer resp.Body.Close()
	c <- Reply{
		time.Since(start),
		resp.StatusCode,
		nil,
	}
	//fmt.Println(resp.Status) OMIT
	//io.Copy(os.Stdout, resp.Body) OMIT
	<-w
}
