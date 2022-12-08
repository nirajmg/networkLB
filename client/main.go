package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup
var client *http.Client
var requestURL = "http://localhost:9000/get"
var csvFile = "csv/demo.csv"
var ips [][]string

func sendRequest(m *sync.Mutex) {
	jsonBody := []byte(`{"UserId":"niraj","Description":"test"}`)
	bodyReader := bytes.NewReader(jsonBody)

	start := time.Now()
	req, err := http.NewRequest(http.MethodGet, requestURL, bodyReader)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		wg.Done()
		panic(err)
	}
	defer resp.Body.Close()
	endTime := time.Now().Sub(start).String()
	// fmt.Println(resp.StatusCode)
	m.Lock()
	ips = append(ips, []string{resp.Header.Get("X-Metrics-IP"), endTime, resp.Header.Get("X-Metrics-Time")})
	m.Unlock()
	wg.Done()

}

func worker(id int, jobs <-chan int) {
	var m sync.Mutex
	for _ = range jobs {
		sendRequest(&m)
	}
}

func main() {
	const numJobs = 100
	jobs := make(chan int, numJobs)

	wg.Add(numJobs)

	for w := 1; w <= 100; w++ {
		go worker(w, jobs)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	wg.Wait()
	close(jobs)

	fmt.Println(len(ips))
	csvFile, err := os.Create(csvFile)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, ip := range ips {
		_ = csvwriter.Write(ip)
	}
	csvwriter.Flush()
	csvFile.Close()

}
