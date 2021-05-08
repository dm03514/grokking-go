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

func TestLogicalRace(t *testing.T) {
	flag.Parse()

	reqCount := SynchronizedCounter{
		mu: &sync.Mutex{},
	}

	go func() {
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := reqCount.Value()
			fmt.Printf("handling request: %d\n", value)
			time.Sleep(1 * time.Nanosecond)
			reqCount.Set(value + 1)
			fmt.Fprintln(w, "Hello, client")
		}))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	requestsChan := make(chan int)

	var wg sync.WaitGroup
	wg.Add(numConcurrentRequests)

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

	fmt.Printf("Num Requests TO Make: %d\n", numRequestsToMake)
	fmt.Printf("Final Count: %d\n", reqCount.Value())
	if numRequestsToMake != reqCount.Value() {
		t.Errorf("expected %d requests: received %d", numRequestsToMake, reqCount.Value())
	}
}
