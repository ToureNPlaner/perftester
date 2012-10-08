package main

import (
	"flag"
	"fmt"
	"github.com/ToureNPlaner/perftester/perftests"
	"net/http"
	"time"
	"math"
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
		test := perftests.NewStdAlgTest(server)
		go doRequests(test, numRequests, resChan)
	}
	var num, failed uint
	var duration, sumTime int64
	var mean, delta, M2 float64

	for i := uint(0); i < numConcurrent*numRequests; i++ {
		res := <-resChan
		num++
		duration = res.Duration.Nanoseconds()
		delta = float64(duration) - mean
		mean += delta/float64(num)
		M2 += delta*(float64(duration)-mean)
		sumTime += duration
		if res.HttpStatus != 200 {
			failed++
			fmt.Printf("Request failed with status %d\n", res.HttpStatus)
		}
	}

	wallSumTime := time.Duration(int64(sumTime) / int64(numConcurrent))
	//variance_n := M2 / float64(num)     // Population Variance
	variance := M2 / (float64(numRequests - 1)) // Sample Variance

	throughput := float64(num) / wallSumTime.Seconds() // in #Reqs/s
	if outputFormat == "human" {
		fmt.Printf("Sent %d Requests (%d failed)\n", num, failed)
		fmt.Printf("Test took: %s (wall time calculated from wait times)\n", time.Duration(wallSumTime))
		fmt.Printf("Average duration: %s\n", time.Duration(int64(mean)))
		fmt.Printf("Standard Deviation: %s\n", time.Duration(int64(math.Sqrt(variance))))
		fmt.Printf("Throughput is: %f #Reqs/s \n", throughput)
	} else if outputFormat == "csv" {
		fmt.Printf("%f, %f, %f\n", float64(math.Sqrt(variance)/float64(time.Millisecond)), float64(mean)/float64(time.Millisecond), throughput)
	}
}
