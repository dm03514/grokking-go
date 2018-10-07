package races

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestDesignNoRace(t *testing.T) {
	flag.Parse()

	reqCount := Counter{}

	var wg sync.WaitGroup
	wg.Add(numConcurrentRequests)
	countChan := make(chan struct{})
	go func() {
		for range countChan {
			reqCount.Inc()
			fmt.Printf("handling request: %d\n", reqCount.Value())
		}
	}()

	go func() {
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Nanosecond)
			fmt.Fprintln(w, "Hello, client")
			countChan <- struct{}{}
		}))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	requestsChan := make(chan int)

	// start a pool of 100 workers all making requests
	for i := 0; i < numConcurrentRequests; i++ {
		go func() {
			defer wg.Done()
			for range requestsChan {
				res, err := http.Get("http://localhost:8080/")
				if err != nil {
					t.Fatal(err)
				}
				_, err = ioutil.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}
		}()
	}

	go func() {
		for i := 0; i < numRequestsToMake; i++ {
			requestsChan <- i
		}
		close(requestsChan)
	}()

	wg.Wait()
	close(countChan)

	fmt.Printf("Num Requests TO Make: %d\n", numRequestsToMake)
	fmt.Printf("Final Count: %d\n", reqCount.Value())
	if numRequestsToMake != reqCount.Value() {
		t.Errorf("expected %d requests: received %d", numRequestsToMake, reqCount.Value())
	}
}
