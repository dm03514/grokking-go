package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCounterMonitor(t *testing.T) {
	var wg = sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	in := NewCounterMonitor(ctx)
	in <- 1
	close(in)
}

func TestSafeCounter(t *testing.T) {
	sc := SafeCounter{
		mu: &sync.Mutex{},
	}
	for i := 0; i < 10; i++ {
		sc.Inc()
	}
	if sc.count != 10 {
		t.Errorf("expecting count to be 10, received: %d", sc.count)
	}
}

func benchmarkCounterMonitor(numProducingGoroutines int, b *testing.B) {
	var wg = sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := NewCounterMonitor(ctx)
	fmt.Printf("Num producing goroutines: %d\n", numProducingGoroutines)
	if numProducingGoroutines == 0 {
		for n := 0; n < b.N; n++ {
			in <- 1
		}
		return
	}

	// instantiate the correct # of senders
	chs := []chan struct{}{}
	for i := 0; i < numProducingGoroutines; i++ {
		wg.Add(1)
		triggerSend := make(chan struct{})
		NewTestSender(
			ctx,
			&wg,
			in,
			triggerSend,
		)
		chs = append(chs, triggerSend)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	// we want to only do N writes across all go routines
	// ie trying to test contention amount routines
	for n := 0; n < b.N; n++ {
		// get a random
		toSend := chs[r1.Intn(len(chs))]
		toSend <- struct{}{}
		/*
		for _, toSend := range chs {
			toSend <- struct{}{}
		}
		*/
	}

	for _, toSend := range chs {
		close(toSend)
	}
	wg.Wait()
	close(in)
}


func BenchmarkCounterMonitorFromMainGoroutine(b *testing.B) {
	benchmarkCounterMonitor(0, b)
}

func BenchmarkCounterMonitor1Goroutine(b *testing.B) {
	benchmarkCounterMonitor(1, b)
}

func BenchmarkCounterMonitor10Goroutines(b *testing.B) {
	benchmarkCounterMonitor(10, b)
}

func BenchmarkCounterMonitor100Goroutines(b *testing.B) {
	benchmarkCounterMonitor(100, b)
}

func BenchmarkCounterMonitor1000Goroutines(b *testing.B) {
	benchmarkCounterMonitor(1000, b)
}

func BenchmarkCounterMonitor10000Goroutines(b *testing.B) {
	benchmarkCounterMonitor(10000, b)
}

func benchmarkSafeCounter(numProducingGoroutines int, b *testing.B) {
	sc := SafeCounter{
		mu: &sync.Mutex{},
	}
	if numProducingGoroutines == 0 {
		for n := 0; n < b.N; n++ {
			sc.Inc()
		}
	}
}

func BenchmarkSafeCounterFromMainGoroutine(b *testing.B) {
	benchmarkSafeCounter(0, b)
}
