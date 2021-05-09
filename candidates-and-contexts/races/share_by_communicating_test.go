package races

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
)

type CounterMonitor struct {
	count int
}

func (cm *CounterMonitor) Monitor(c chan struct{}) {
	for range c {
		cm.count++
	}
}

func (cm *CounterMonitor) Count() int {
	return cm.count
}

func TestMonitorNoRace(t *testing.T) {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(numConcurrentRequests)

	countChan := make(chan struct{})

	cm := &CounterMonitor{}
	go cm.Monitor(countChan)

	go func() {
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("Final Count: %d\n", cm.Count())
	if numRequestsToMake != cm.Count() {
		t.Errorf("expected %d requests: received %d", numRequestsToMake, cm.Count())
	}
}
