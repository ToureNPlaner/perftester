package perftests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type StdAlgTest struct {
	Server    string
}

var upperLat, lowerLat, leftLon, rightLon float64
var numPoints uint
var algSuffix string

func init() {
	flag.StringVar(&algSuffix, "algorithm", "sp", "the algorithm suffix to use")
	flag.Float64Var(&upperLat, "upperLat", 54.0, "upper latitude")
	flag.Float64Var(&lowerLat, "lowerLat", 47.0, "lower latitude")
	flag.Float64Var(&leftLon, "leftLon", 5.9, "left longitude")
	flag.Float64Var(&rightLon, "rightLon", 14.9, "right longitude")
	flag.UintVar(&numPoints, "numPoints", 2, "number of points in request")
	rand.Seed(42); // So we get always the same points
}

type tpPoint struct {
	Lt int32 `json:"lt"`
	Ln int32 `json:"ln"`
}

type tpRequest struct {
	Points []tpPoint `json:"points"`
}

func NewStdAlgTest(server string) (res *StdAlgTest) {
	res = new(StdAlgTest)
	res.Server = server
	return
}

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func (r *StdAlgTest) DoRequest(client *http.Client, resChan chan PerfResult) {
	var tpReq tpRequest
	tpReq.Points = make([]tpPoint, numPoints)
	for i := uint(0); i < numPoints; i++ {
		tpReq.Points[int(i)] = tpPoint{int32(random(lowerLat, upperLat) * 1.0e+7), int32(random(leftLon, rightLon) * 1.0e+7)}
	}

	b, err := json.Marshal(tpReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	startTime := time.Now()
	response, err := client.Post(r.Server+"/alg"+algSuffix, "application/x-jackson-smile", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
		return
	}
	response.Body.Close()
	resChan <- PerfResult{response.StatusCode, time.Since(startTime)}
}
