package perftests

import (
   "net/http"
   "encoding/json"
   "bytes"
   "fmt"
   "time"
)

type StdAlgTest struct {
   Server string
   AlgSuffix string
}

type tpPoint struct {
   Lt uint32 `json:"lt"`
   Ln uint32 `json:"ln"`
}

type tpRequest struct {
   Points []tpPoint `json:"points"`
}

func (r *StdAlgTest) DoRequest(client *http.Client, resChan chan PerfResult){
   tpReq := tpRequest{[]tpPoint{tpPoint{Lt: 487786110, Ln: 91794440}, tpPoint{Lt: 535652780, Ln: 100013890}}}
   b, err := json.Marshal(tpReq)
   if err != nil {
      fmt.Println(err)
      return
   }
   startTime := time.Now()
   response, err := client.Post(r.Server+"alg"+r.AlgSuffix,"application/json", bytes.NewBuffer(b))
   if err != nil {
      fmt.Println(err)
      return
   }
   resChan <- PerfResult{response.StatusCode, time.Since(startTime)}
}
