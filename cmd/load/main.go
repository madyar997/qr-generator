package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	endpoints = []string{
		"/code-200",
		"/code-400",
		"/code-500",
	}
	requests = map[string]int{}
)

func randomizeEndpoints() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(endpoints),
		func(i, j int) { endpoints[i], endpoints[j] = endpoints[j], endpoints[i] },
	)
}

func main() {
	maxSuccessfulRequests, maxErrorRequests := 500, 200

	totalRequestsCount := 0

	randomizeEndpoints()
	for {
		for _, e := range endpoints {
			requestsToEndpoint := 0
			if e == "/code-200" {
				requestsToEndpoint = rand.Intn(maxSuccessfulRequests)
			} else {
				requestsToEndpoint = rand.Intn(maxErrorRequests)
			}

			requests[e] = requestsToEndpoint
			totalRequestsCount += requestsToEndpoint
		}

		log.Println(totalRequestsCount)

		wg := &sync.WaitGroup{}
		wg.Add(totalRequestsCount)

		for endpoint, requestsCount := range requests {
			for i := 0; i < requestsCount; i++ {
				go func(e string) {
					if _, err := http.DefaultClient.Get(fmt.Sprintf("http://localhost:8081%s", e)); err != nil {
						fmt.Printf("is it err: %v\n", err)
					}
					wg.Done()
				}(endpoint)
			}
		}

		wg.Wait()
	}
}
