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
	flag.StringVar(&outputFormat, "format", "human", "The format to use, choose from human,csv (average time in ms, average request size, throughput, data throughput)")
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
	var duration, sumTime, size, sumBytes int64
	var meanDuration, deltaDuration, m2Duration float64
	var meanSize, deltaSize, m2Size float64

	for i := uint(0); i < numConcurrent*numRequests; i++ {
		res := <-resChan
		num++
		duration = res.Duration.Nanoseconds()
		size = res.ContentLength

		deltaDuration = float64(duration) - meanDuration
		deltaSize = float64(size) - meanSize

		meanDuration += deltaDuration/float64(num)
		meanSize += deltaSize/float64(num)

		m2Duration += deltaDuration*(float64(duration)-meanDuration)
		m2Size += deltaSize*(float64(size)-meanSize)

		sumTime += duration
		sumBytes += size
		if res.HttpStatus != 200 {
			failed++
			fmt.Printf("Request failed with status %d\n", res.HttpStatus)
		}
	}

	wallSumTime := time.Duration(int64(sumTime) / int64(numConcurrent))
	//variance_n := M2 / float64(num)     // Population Variance
	varianceDuration := m2Duration / float64(numRequests - 1) // Sample Variance
	varianceSize := m2Size / float64(numRequests - 1)

	throughput := float64(num) / wallSumTime.Seconds() // in #Reqs/s
	networkThroughput := (float64(sumBytes*8)/(1000000))/ wallSumTime.Seconds() // in MBit/s

	if outputFormat == "human" {
		fmt.Printf("Sent %d Requests (%d failed)\n", num, failed)
		fmt.Printf("Test took: %s (wall time calculated from wait times)\n", time.Duration(wallSumTime))
		fmt.Printf("Throughput is: %f #Reqs/s \n\n", throughput)
		fmt.Printf("Average duration: %s\n", time.Duration(int64(meanDuration)))
		fmt.Printf("Standard Deviation: %s\n\n", time.Duration(int64(math.Sqrt(varianceDuration))))
		fmt.Printf("Data transfered: %f MB\n", float64(sumBytes)/(1024.0*1024.0))
		fmt.Printf("Average Resultsize: %f KB\n", meanSize/1024.0)
		fmt.Printf("Standard Deviation: %f KB\n", math.Sqrt(varianceSize)/1024.0)
		fmt.Printf("Network Throughput is: %f MBit/s\n", networkThroughput)
	} else if outputFormat == "csv" {
		fmt.Printf("%f, %f, %f, %f\n", float64(int64(meanDuration))/float64(time.Millisecond), meanSize/1024.0, throughput, networkThroughput)
	}
}
