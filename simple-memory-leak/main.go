package main

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

type RequestTracker struct {
	mu       sync.Mutex
	requests [][]byte
}

func (rt *RequestTracker) Track(req *http.Request) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	// alloc 10KB for each track
	rt.requests = append(rt.requests, bytes.Repeat([]byte("a"), 10000))
}

var (
	requests RequestTracker

	responseLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "response_latency_seconds",
		Help: "Length of HTTP processing.",
	})
)

func init() {
	requests = RequestTracker{}
	prometheus.MustRegister(responseLatency)
}

func main() {
	fmt.Printf("Starting Server on port :%s\n", "8080")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requests.Track(r)
		r.Body.Close()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))

		responseLatency.Observe((time.Since(start)).Seconds())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
