package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup
var client *http.Client

var ips []string

func sendRequest(m *sync.Mutex) {

	resp, err := client.Get("http://localhost:30000")
	if err != nil {
		panic(err)
	}

	m.Lock()
	ips = append(ips, resp.Header.Get("X-Metrics-IP"))
	m.Unlock()
	wg.Done()
}

func worker(id int, jobs <-chan int) {
	var m sync.Mutex
	for _ = range jobs {
		go sendRequest(&m)
	}
}

func main() {

	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100
	client = &http.Client{Transport: &defaultTransport}

	const numJobs = 100
	jobs := make(chan int, numJobs)

	wg.Add(100)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 100; w++ {
		go worker(w, jobs)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	wg.Wait()
	fmt.Println(ips)
	csvFile, err := os.Create("iphash.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, ip := range ips {
		_ = csvwriter.Write([]string{ip})
	}
	csvwriter.Flush()
	csvFile.Close()
	close(jobs)
}
