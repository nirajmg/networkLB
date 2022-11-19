package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup
var client *http.Client

func sendRequest(filename string, m *sync.Mutex) {

	resp, err := client.Get("http://localhost:30924")
	if err != nil {
		panic(err)
	}

	// ip := resp.Header.Get("IP")

	// fmt.Println(resp.Header.Get("Date"))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	b := []string{string(body)}

	m.Lock()
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	// defer w.Flush()
	w.Write(b)
	m.Unlock()
	wg.Done()
}

// func main() {

// 	var m sync.Mutex

// 	f, _ := os.Create("roundrobin.csv")
// 	defer f.Close()

// 	for j := 0; j < 10; j++ {

// 	}
// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go sendRequest("roundrobin.csv", &m, &wg)
// 	}

// 	wg.Wait()
// }

func worker(id int, jobs <-chan int) {
	var m sync.Mutex
	for _ = range jobs {
		go sendRequest("roundrobin.csv", &m)
	}
}

func main() {

	tr := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client = &http.Client{Transport: tr}
	// resp, err := client.Get("https://example.com")

	f, _ := os.Create("roundrobin.csv")
	defer f.Close()

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	const numJobs = 1000
	jobs := make(chan int, numJobs)

	wg.Add(1000)

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
	close(jobs)
}
