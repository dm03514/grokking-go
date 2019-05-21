package main

import (
	"context"
	"fmt"
	"sync"
)

func NewCounterMonitor(ctx context.Context) chan<- int {
	counter := 0
	ch := make(chan int)

	go func() {
		for {
			select {
			case i, ok := <-ch:
				if !ok {
					fmt.Printf("NewCounterMonitor: final_count: %d\n", counter)
					return
				}
				counter += i
			case <-ctx.Done():
				fmt.Printf("context Cancelled: final_count: %d\n", counter)
				return
			}
		}
	}()

	return ch
}

type SafeCounter struct {
	mu    *sync.Mutex
	count int
}

func (s *SafeCounter) Inc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
}


func NewTestSender(ctx context.Context, wg *sync.WaitGroup, counterCh chan<- int, triggerSend <-chan struct{}) {
	go func() {
		defer wg.Done()

		for {
			select {
				case _, ok := <-triggerSend:
					if !ok {
						return
					}

					counterCh <- 1

				case <-ctx.Done():
					fmt.Printf("new sender context cancelled")
					return
			}
		}
	}()
}
