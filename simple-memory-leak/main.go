package main

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"sync"
)

type RequestTracker struct {
	mu sync.Mutex
	requests []*http.Request
}
func (rt *RequestTracker) Track(req *http.Request) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.requests = append(rt.requests, req)
}

var requests RequestTracker

func init() {
	requests = RequestTracker{}
}


func main() {
	fmt.Print("hello")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requests.Track(r)
		r.Body.Close()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}