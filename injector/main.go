package main

import "github.com/owulveryck/api-repository/object"

func main() {
	work := make(chan object.IDer)

	// startup pool of 10 go routines and read urls from work channel
	for i := 0; i <= 10; i++ {
		go func(w chan object.IDer) {
			url := <-w
		}(work)
	}

	// write urls to the work channel, blocking until a worker goroutine
	// is able to start work
	for _, url := range urlsToProcess {
		work <- url
	}
}
