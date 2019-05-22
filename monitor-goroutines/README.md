


## Benchmarks

```
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/dm03514/grokking-go/monitor-goroutines
BenchmarkCounterMonitorFromMainGoroutine-12              5000000               292 ns/op
BenchmarkCounterMonitor1Goroutine-12                     1000000              1022 ns/op
BenchmarkCounterMonitor10Goroutines-12                   1000000              1101 ns/op
BenchmarkCounterMonitor100Goroutines-12                  1000000              1232 ns/op
BenchmarkCounterMonitor1000Goroutines-12                 1000000              1565 ns/op
BenchmarkCounterMonitor10000Goroutines-12                1000000              1426 ns/op
BenchmarkCounterMonitorFromMainGoroutineBuffered-12      5000000               261 ns/op
BenchmarkCounterMonitor1GoroutineBuffered-12             1000000              1081 ns/op
BenchmarkCounterMonitor10GoroutinesBuffered-12           2000000               931 ns/op
BenchmarkCounterMonitor100GoroutinesBuffered-12          1000000              1364 ns/op
BenchmarkCounterMonitor1000GoroutinesBuffered-12         1000000              1237 ns/op
BenchmarkCounterMonitor10000GoroutinesBuffered-12        1000000              1420 ns/op
BenchmarkSafeCounterFromMainGoroutine-12                30000000                41.8 ns/op
BenchmarkSafeCounter1Goroutine-12                        3000000               432 ns/op
BenchmarkSafeCounter10Goroutine-12                       2000000               672 ns/op
BenchmarkSafeCounter100Goroutine-12                      2000000               908 ns/op
BenchmarkSafeCounter1000Goroutine-12                     1000000              1007 ns/op
BenchmarkSafeCounter10000Goroutine-12                    2000000              1199 ns/op
PASS
ok      github.com/dm03514/grokking-go/monitor-goroutines       30.996s
```