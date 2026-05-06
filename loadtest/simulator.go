package loadtest

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Transaction struct {
	ID     int
	Amount int64
}

type Result struct {
	Total         int
	Processed     int
	Success       int
	Failed        int
	TotalAmount   int64
	AverageLatency time.Duration
	Duration      time.Duration
	ThroughputTPS float64
}

type workerResult struct {
	success bool
	amount  int64
	latency time.Duration
}

// SimulateConcurrentTransactions runs a worker-pool-based simulation for transactions.
func SimulateConcurrentTransactions(totalTransactions int, workerCount int) (Result, error) {
	if totalTransactions <= 0 {
		return Result{}, errors.New("totalTransactions must be > 0")
	}
	if workerCount <= 0 {
		return Result{}, errors.New("workerCount must be > 0")
	}

	jobs := make(chan Transaction, totalTransactions)
	results := make(chan workerResult, totalTransactions)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	start := time.Now()
	for i := 0; i < totalTransactions; i++ {
		jobs <- Transaction{
			ID:     i + 1,
			Amount: int64(100 + (i % 1000)),
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var mu sync.Mutex
	aggregated := Result{Total: totalTransactions}
	var totalLatency time.Duration

	for r := range results {
		mu.Lock()
		aggregated.Processed++
		if r.success {
			aggregated.Success++
			aggregated.TotalAmount += r.amount
		} else {
			aggregated.Failed++
		}
		totalLatency += r.latency
		mu.Unlock()
	}

	aggregated.Duration = time.Since(start)
	if aggregated.Processed > 0 {
		aggregated.AverageLatency = totalLatency / time.Duration(aggregated.Processed)
	}
	if aggregated.Duration > 0 {
		aggregated.ThroughputTPS = float64(aggregated.Processed) / aggregated.Duration.Seconds()
	}

	return aggregated, nil
}

func worker(_ int, jobs <-chan Transaction, results chan<- workerResult, wg *sync.WaitGroup) {
	defer wg.Done()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for tx := range jobs {
		latency := time.Duration(1+rng.Intn(5)) * time.Millisecond
		time.Sleep(latency)

		// Simulate a small failure ratio to imitate realistic processing.
		success := rng.Intn(100) >= 2
		results <- workerResult{
			success: success,
			amount:  tx.Amount,
			latency: latency,
		}
	}
}
