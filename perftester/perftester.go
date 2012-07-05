package main

import (
	"flag"
	"fmt"
	"net/http"
	"org.tourenplaner/perftests"
)

var server, algSuffix string
var numRequests, numConcurrent uint

func init() {
	flag.StringVar(&server, "server", "http://localhost:8080/", "The server to connect to (as URL)")
	flag.StringVar(&algSuffix, "alg", "sp", "the algorithm suffix to use")
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
	fmt.Printf("Connecting to server %s requesting algorithm for %s \n", server, algSuffix)
	resChan := make(chan perftests.PerfResult, numConcurrent)

	for i := uint(0); i < numConcurrent; i++ {
		test := &perftests.StdAlgTest{server, algSuffix}
		go doRequests(test, numRequests, resChan)
	}

	for i := uint(0); i < numConcurrent*numRequests; i++ {
		res := <-resChan
		fmt.Printf("Status: %v Duration: %s\n", res.HttpStatus, res.Duration)
	}
}
