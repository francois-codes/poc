package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	totalRequests = 50000
	numWorkers    = 10
	minID         = 33
	maxID         = 1032
	baseURL       = "http://127.0.0.1:8585/datamodel/%d"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Channel to distribute work
	jobs := make(chan int, totalRequests)
	// Channel to collect durations
	results := make(chan time.Duration, totalRequests)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Send jobs
	fmt.Println("Starting the queries")
	start := time.Now()
	for i := 0; i < totalRequests; i++ {
		jobs <- rand.Intn(maxID-minID+1) + minID
	}
	close(jobs)

	// Wait for workers
	wg.Wait()
	end := time.Since(start).Seconds()
	close(results)

	// Collect and analyze timings
	var durations []float64
	for d := range results {
		durations = append(durations, float64(d.Microseconds()))
	}

	if len(durations) == 0 {
		log.Println("No durations collected.")
		return
	}

	qps := float64(len(durations)) / end

	sort.Float64s(durations)

	// Compute stats
	min := durations[0]
	max := durations[len(durations)-1]
	sum := 0.0
	for _, d := range durations {
		sum += d
	}
	avg := sum / float64(len(durations))
	p95 := percentile(durations, 95)

	// Print
	fmt.Println("📊 HTTP Benchmark Stats (in microseconds):")
	fmt.Printf("   ➤ Fastest: %.0f µs\n", min)
	fmt.Printf("   ➤ Slowest: %.0f µs\n", max)
	fmt.Printf("   ➤ Average: %.0f µs\n", avg)
	fmt.Printf("   ➤ P95:     %.0f µs\n", p95)
	fmt.Printf("⚡ Total time: %.2fs\n", end)
	fmt.Printf("🚀 QPS (Queries per second): %.2f\n", qps)
}

func worker(jobs <-chan int, results chan<- time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	for id := range jobs {
		url := fmt.Sprintf(baseURL, id)
		start := time.Now()
		resp, err := client.Get(url)
		elapsed := time.Since(start)
		if err == nil && resp.Body != nil {
			resp.Body.Close()
		}
		results <- elapsed
	}
}

func percentile(sorted []float64, percent float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := (percent / 100.0) * float64(len(sorted)-1)
	i := int(math.Floor(index))
	if i == len(sorted)-1 {
		return sorted[i]
	}
	frac := index - float64(i)
	return sorted[i]*(1-frac) + sorted[i+1]*frac
}

/*
📊 HTTP Benchmark Stats :
   ➤ Fastest: 0.255 ms
   ➤ Slowest: 77.253 ms
   ➤ Average: 3.873 ms
   ➤ P95:     8.004 ms
⚡ Total time: 19.39s
🚀 QPS (Queries per second): 2578.77
*/
