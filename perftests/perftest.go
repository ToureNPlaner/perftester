package perftests

import (
        "time"
        "net/http"
)


type PerfResult struct {
        HttpStatus int
        Duration time.Duration
}

type PerfTest interface{
        DoRequest(cl *http.Client, results chan PerfResult)
}
