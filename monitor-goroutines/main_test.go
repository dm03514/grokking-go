package main

import (
	"context"
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
	in := NewCounterMonitor(ctx, &wg,0)
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

func benchmarkCounterMonitor(numProducingGoroutines int, channelBufferSize int, b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var counterMonitorWG = sync.WaitGroup{}
	counterMonitorWG.Add(1)
	in := NewCounterMonitor(ctx, &counterMonitorWG, channelBufferSize)
	// fmt.Printf("Num producing goroutines: %d\n", numProducingGoroutines)
	if numProducingGoroutines == 0 {
		for n := 0; n < b.N; n++ {
			in <- 1
		}
		return
	}

	// instantiate the correct # of senders
	chs := []chan struct{}{}
	var wg = sync.WaitGroup{}
	for i := 0; i < numProducingGoroutines; i++ {
		wg.Add(1)
		triggerSend := make(chan struct{})
		NewTestSender(
			ctx,
			&wg,
			func() {
				in <- 1
			},
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
	}

	for _, toSend := range chs {
		close(toSend)
	}

	wg.Wait()
	close(in)
	counterMonitorWG.Wait()
}

func BenchmarkCounterMonitorFromMainGoroutine(b *testing.B) {
	benchmarkCounterMonitor(0, 0, b)
}

func BenchmarkCounterMonitor1Goroutine(b *testing.B) {
	benchmarkCounterMonitor(1, 0, b)
}

func BenchmarkCounterMonitor10Goroutines(b *testing.B) {
	benchmarkCounterMonitor(10, 0, b)
}

func BenchmarkCounterMonitor100Goroutines(b *testing.B) {
	benchmarkCounterMonitor(100, 0, b)
}

func BenchmarkCounterMonitor1000Goroutines(b *testing.B) {
	benchmarkCounterMonitor(1000, 0, b)
}

func BenchmarkCounterMonitor10000Goroutines(b *testing.B) {
	benchmarkCounterMonitor(10000, 0, b)
}

func benchmarkSafeCounter(numProducingGoroutines int, b *testing.B) {
	var wg = sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc := SafeCounter{
		mu: &sync.Mutex{},
	}
	if numProducingGoroutines == 0 {
		for n := 0; n < b.N; n++ {
			sc.Inc()
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
			func() {
				sc.Inc()
			},
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
	}

	for _, toSend := range chs {
		close(toSend)
	}
	wg.Wait()
}

func BenchmarkSafeCounterFromMainGoroutine(b *testing.B) {
	benchmarkSafeCounter(0, b)
}

func BenchmarkSafeCounter1Goroutine(b *testing.B) {
	benchmarkSafeCounter(1, b)
}

func BenchmarkSafeCounter10Goroutine(b *testing.B) {
	benchmarkSafeCounter(10, b)
}

func BenchmarkSafeCounter100Goroutine(b *testing.B) {
	benchmarkSafeCounter(100, b)
}

func BenchmarkSafeCounter1000Goroutine(b *testing.B) {
	benchmarkSafeCounter(1000, b)
}

func BenchmarkSafeCounter10000Goroutine(b *testing.B) {
	benchmarkSafeCounter(10000, b)
}
