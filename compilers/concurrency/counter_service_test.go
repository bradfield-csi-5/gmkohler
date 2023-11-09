// Author: Patch Neranartkomol

package counterservice

import (
	"runtime"
	"sync"
	"testing"
)

const (
	NumGoroutines     = 32
	CallsPerGoroutine = 100000
)

func getNextMonotonicityChecker(counter CounterService, t *testing.T) {
	var prev uint64
	for i := 0; i < CallsPerGoroutine; i++ {
		value := counter.GetNext()
		if value <= prev {
			t.Fatalf("Values were NOT monotonically increasing; value: %d, prev: %d", value, prev)
		}
		prev = value
	}
}

func TestSynchronizedCounterServices(t *testing.T) {
	var wg sync.WaitGroup
	ccs := NewChannelCounterService()
	defer ccs.Close()
	counters := []CounterService{
		&AtomicCounterService{},
		&MutexCounterService{},
		ccs,
	}
	for _, counter := range counters {
		for i := 0; i < NumGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				getNextMonotonicityChecker(counter, t)
			}()
		}
		wg.Wait()
		nextVal := counter.GetNext()
		if nextVal != (NumGoroutines*CallsPerGoroutine)+1 {
			t.Errorf("Counter ID does not match total calls; nextVal: %d", nextVal)
		}
	}
}

// This test only checks that the unsynchronized version is correct when run in a single goroutine.
// It does not spawn additional goroutines.
func TestUnsynchronizedCounterService(t *testing.T) {
	var counter CounterService = &UnsynchronizedCounterService{}
	getNextMonotonicityChecker(counter, t)
	nextVal := counter.GetNext()
	if nextVal != CallsPerGoroutine+1 {
		t.Fatalf("Counter ID does not match total calls; nextVal: %d", nextVal)
	}
}

func BenchmarkCounterServices(b *testing.B) {
	ccs := NewChannelCounterService()
	defer ccs.Close()
	for _, testCase := range []struct {
		name    string
		counter CounterService
	}{
		{"unsynchronized", &UnsynchronizedCounterService{}},
		{"atomic", &AtomicCounterService{}},
		{"mutex", &MutexCounterService{}},
		{"channel", ccs},
	} {
		b.Run(testCase.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				testCase.counter.GetNext()
			}
		})
	}
}

func BenchmarkCounterServicesContended(b *testing.B) {
	b.SetParallelism(NumGoroutines / runtime.GOMAXPROCS(0))
	ccs := NewChannelCounterService()
	defer ccs.Close()
	for _, testCase := range []struct {
		name    string
		counter CounterService
	}{
		{"unsynchronized", &UnsynchronizedCounterService{}},
		{"atomic", &AtomicCounterService{}},
		{"mutex", &MutexCounterService{}},
		{"channel", ccs},
	} {
		b.Run(testCase.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					testCase.counter.GetNext()
				}
			})
		})
	}
}
