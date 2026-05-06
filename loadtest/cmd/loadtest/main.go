package main

import (
	"fmt"
	"log"

	"loadtest"
)

func main() {
	totalTransactions := 1000
	workerCount := 100

	result, err := loadtest.SimulateConcurrentTransactions(totalTransactions, workerCount)
	if err != nil {
		log.Fatalf("simulation failed: %v", err)
	}

	fmt.Println("=== Load Test Report ===")
	fmt.Printf("Total Transactions: %d\n", result.Total)
	fmt.Printf("Processed: %d\n", result.Processed)
	fmt.Printf("Success: %d\n", result.Success)
	fmt.Printf("Failed: %d\n", result.Failed)
	fmt.Printf("Total Amount (success): %d\n", result.TotalAmount)
	fmt.Printf("Average Latency: %s\n", result.AverageLatency)
	fmt.Printf("Total Duration: %s\n", result.Duration)
	fmt.Printf("Throughput (TPS): %.2f\n", result.ThroughputTPS)
}
