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
var outputFormat string

func init() {
	flag.StringVar(&outputFormat, "format", "human", "The format to use, choose from human,csv (average time in ms, throughput)")
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
	fmt.Errorf("Connecting to server %s \n", server)
	resChan := make(chan perftests.PerfResult, numConcurrent)

	for i := uint(0); i < numConcurrent; i++ {
		test :=  perftests.NewStdAlgTest(server)
		go doRequests(test, numRequests, resChan)
	}
	sumTime := int64(0)
	failed :=int64(0)
	for i := uint(0); i < numConcurrent*numRequests; i++ {
		res := <-resChan
		sumTime += res.Duration.Nanoseconds()
		if res.HttpStatus != 200 {
			failed++
			fmt.Printf("Request failed with status %d\n", res.HttpStatus)
		}
	}
	totalRequests := int64(numConcurrent*numRequests)
	wallSumTime := int64(sumTime)/int64(numConcurrent)

	averageTime := sumTime/totalRequests // in nanoseconds
	throughput := float64(totalRequests)/time.Duration(wallSumTime).Seconds() // in #Reqs/s
	if outputFormat == "human" {
		fmt.Printf("Sent %d Requests (%d failed)\n", totalRequests, failed)
		fmt.Printf("Test took: %s (wall time calculated from wait times)\n", time.Duration(wallSumTime))
		fmt.Printf("Average duration: %s\n", time.Duration(averageTime))
		fmt.Printf("Throughput is: %f #Reqs/s \n", throughput)
	} else if outputFormat  == "csv" {
		fmt.Printf("%f, %f\n", float64(time.Duration(averageTime))/float64(time.Millisecond), throughput)
	}
}
