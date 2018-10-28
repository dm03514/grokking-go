package main

import (
	"net/http"
	"time"
	"fmt"
	"log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
	requestLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_seconds",
		Help: "Distribution of request lengths",
	})
)

func init() {
	prometheus.MustRegister(requestLatency)
}

type Handler struct {
	AdditionalLatency time.Duration
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	time.Sleep(h.AdditionalLatency)
	/*
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		msg := fmt.Sprintf("received: %q\n", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	*/

	diff := time.Since(start)
	fmt.Fprintf(w, "Took: %s\n", diff)
	requestLatency.Observe(diff.Seconds())
}

func main() {
	fmt.Printf("starting_server: :8080\n")
	http.Handle("/slow", Handler{
		AdditionalLatency: time.Millisecond * 500,
	})
	http.Handle("/fast", Handler{
		AdditionalLatency: time.Millisecond * 0,
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}