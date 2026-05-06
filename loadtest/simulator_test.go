package loadtest

import "testing"

func TestSimulateConcurrentTransactions(t *testing.T) {
	result, err := SimulateConcurrentTransactions(1000, 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Processed != 1000 {
		t.Fatalf("expected processed=1000, got %d", result.Processed)
	}
	if result.Success+result.Failed != 1000 {
		t.Fatalf("expected success+failed=1000, got %d", result.Success+result.Failed)
	}
}

func Benchmark1000ConcurrentTransactions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := SimulateConcurrentTransactions(1000, 100)
		if err != nil {
			b.Fatalf("simulation error: %v", err)
		}
	}
}
