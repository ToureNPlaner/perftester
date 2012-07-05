package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"github.com/ToureNPlaner/perftester/perftests"
)

var server string
var numRequests, numConcurrent uint

func init() {
	flag.StringVar(&server, "server", "http://localhost:8080", "The server to connect to (as URL)")
	flag.UintVar(&numRequests, "requests", 1000, "the number of requests to send")
	flag.UintVar(&numConcurrent, "concurrent", 16, "the number of concurrent connections")
}

func doRequests(test perftests.PerfTest, numRequests uint, resChan chan perftests.PerfResult) {
	var client http.Client
	for i := uint(0); i < numRequests; i++ {
		test.DoRequest(&client, resChan)
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Connecting to server %s \n", server)
	resChan := make(chan perftests.PerfResult, numConcurrent)

	for i := uint(0); i < numConcurrent; i++ {
		test :=  perftests.NewStdAlgTest(server)
		go doRequests(test, numRequests, resChan)
	}
	var sumTime int64 = 0
	for i := uint(0); i < numConcurrent*numRequests; i++ {
		res := <-resChan
		sumTime += res.Duration.Nanoseconds()
		fmt.Printf("Status: %v Duration: %s\n", res.HttpStatus, res.Duration)
	}
	sumTime /= int64(numConcurrent*numRequests)

	fmt.Printf("Average duration: %s\n", time.Duration(sumTime))
}
