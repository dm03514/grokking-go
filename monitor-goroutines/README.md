


## Benchmarks

```
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/dm03514/grokking-go/monitor-goroutines
BenchmarkCounterMonitorFromMainGoroutine-12                      5000000               297 ns/op
BenchmarkCounterMonitor1Goroutine-12                             1000000              1009 ns/op
BenchmarkCounterMonitor10Goroutines-12                           1000000              1100 ns/op
BenchmarkCounterMonitor100Goroutines-12                          1000000              1262 ns/op
BenchmarkCounterMonitor1000Goroutines-12                         1000000              1389 ns/op
BenchmarkCounterMonitor10000Goroutines-12                        1000000              1658 ns/op
BenchmarkContentionCounterMonitor2Goroutines-12                   500000              2681 ns/op
BenchmarkContentionCounterMonitor10Goroutines-12                  100000             12665 ns/op
BenchmarkContentionCounterMonitor100Goroutines-12                  10000            166051 ns/op
BenchmarkContentionCounterMonitor1000Goroutines-12                  1000           1617492 ns/op
BenchmarkContentionCounterMonitor10000Goroutines-12                  100          16229725 ns/op
BenchmarkContentionCounterMonitor2GoroutinesBuffered-12           500000              3071 ns/op
BenchmarkContentionCounterMonitor10GoroutinesBuffered-12          200000             10944 ns/op
BenchmarkContentionCounterMonitor100GoroutinesBuffered-12          10000            149916 ns/op
BenchmarkContentionCounterMonitor1000GoroutinesBuffered-12          1000           1482662 ns/op
BenchmarkContentionCounterMonitor10000GoroutinesBuffered-12          100          16928763 ns/op
BenchmarkSafeCounterFromMainGoroutine-12                        30000000                44.3 ns/op
BenchmarkSafeCounter1Goroutine-12                                3000000               425 ns/op
BenchmarkSafeCounter10Goroutine-12                               2000000               732 ns/op
BenchmarkSafeCounter100Goroutine-12                              2000000               992 ns/op
BenchmarkSafeCounter1000Goroutine-12                             1000000              1255 ns/op
BenchmarkSafeCounter10000Goroutine-12                            1000000              1248 ns/op
BenchmarkContentionSafeCounter2Goroutines-12                     1000000              1690 ns/op
BenchmarkContentionSafeCounter10Goroutine-12                      200000              9074 ns/op
BenchmarkContentionSafeCounter100Goroutine-12                      10000            114061 ns/op
BenchmarkContentionSafeCounter1000Goroutine-12                      1000           1225179 ns/op
BenchmarkContentionSafeCounter10000Goroutine-12                      100          13807604 ns/op
PASS
ok      github.com/dm03514/grokking-go/monitor-goroutines       43.980s
```