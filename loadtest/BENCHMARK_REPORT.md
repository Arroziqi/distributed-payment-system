# Load Test Benchmark Report

## Scenario

- Concurrent transactions: `1000`
- Concurrency model: goroutines + worker pool + channels + mutex-protected aggregation
- Workers: `100`
- Command:

```bash
go test -bench=. -benchmem ./...
```

## Measured Result

```text
Benchmark1000ConcurrentTransactions-8   	      31	  37999668 ns/op	  597676 B/op	     315 allocs/op
```

## Interpretation

- Average benchmark runtime per simulation: `~37.99ms`
- Simulations completed in benchmark: `31`
- Memory per simulation: `~597 KB`
- Allocations per simulation: `315`

## Notes

- The simulator includes a small synthetic failure ratio (~2%) to mimic real payment pipeline behavior under load.
- Latency is simulated between `1-5ms` per transaction worker task.
