package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/owulveryck/api-repository/injector"
)

var concurrency, numElements int

func main() {
	concurrency := flag.Int("c", 1, "concurrency")
	file := flag.String("f", "", "file")
	url := flag.String("u", "", "url")
	flag.Parse()
	if *file == "" || *url == "" {
		log.Fatal("Invalid arguments")
	}
	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	numElements, err := lineCounter(f)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	wg := new(sync.WaitGroup)
	replyChan := make(chan injector.Reply)
	slo := &injector.SLO{
		Latency: map[float64]time.Duration{
			0.95: 400 * time.Millisecond,
			0.99: 750 * time.Millisecond,
		},
		Allowed5xxErrors: 100 - 99.9,
		Verbose:          true,
	}
	wg.Add(1)
	go slo.Evaluate(replyChan, numElements, wg)
	// END_MAIN OMIT
	ws := make(chan struct{}, *concurrency) // Number of concurrent calls // HL
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ws <- struct{}{}
		wg.Add(1)
		go injector.SendPostRequest(*url, scanner.Text(), ws, wg, replyChan) // HL
	}
	wg.Wait()
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
