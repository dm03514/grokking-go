package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"html"
	"log"
	"net/http"
	"time"
)

var (
	requestLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_seconds",
		Help: "Distribution of request lengths",
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(requestLatency)
}

type Handler struct {
	AdditionalLatency time.Duration
	NumBytesToAlloc   int
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	time.Sleep(h.AdditionalLatency)

	bytes.Repeat([]byte("a"), h.NumBytesToAlloc)

	diff := time.Since(start)
	fmt.Fprintf(w, "Took: %s\n", diff)
	requestLatency.WithLabelValues(html.EscapeString(r.URL.Path)).Observe(diff.Seconds())
}

func main() {
	slowDuration := flag.Duration("slow-duration", time.Millisecond*10, "value to delay slow request by")
	numBytesToAlloc := flag.Int("num-bytes-to-alloc", 100, "how many bytes to allocate on each request")
	flag.Parse()

	fmt.Printf("starting_server: :8080\n")
	http.Handle("/slow", Handler{
		AdditionalLatency: *slowDuration,
		NumBytesToAlloc:   *numBytesToAlloc,
	})
	http.Handle("/fast", Handler{
		AdditionalLatency: time.Millisecond * 0,
		NumBytesToAlloc:   0,
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
