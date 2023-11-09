// Author: Patch Neranartkomol

package counterservice

import (
	"sync"
	"sync/atomic"
)

type CounterService interface {
	// GetNext Returns values in ascending order; it should be safe to call
	// GetNext() concurrently from multiple goroutines without any
	// additional synchronization on the caller's side.
	GetNext() uint64
}

type UnsynchronizedCounterService struct {
	counter uint64
}

func (c *UnsynchronizedCounterService) GetNext() uint64 {
	c.counter++
	return c.counter
}

type AtomicCounterService struct {
	counter atomic.Uint64
}

func (c *AtomicCounterService) GetNext() uint64 {
	return c.counter.Add(1)
}

type MutexCounterService struct {
	mu      sync.Mutex
	counter uint64
}

func (c *MutexCounterService) GetNext() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter++
	return c.counter
}

type ChannelCounterService struct {
	counter   uint64
	counterCh chan uint64
	incrCh    chan struct{}
}

// NewChannelCounterService is a constructor for ChannelCounterService
func NewChannelCounterService() *ChannelCounterService {
	cs := ChannelCounterService{
		incrCh:    make(chan struct{}, 20),
		counterCh: make(chan uint64, 20),
	}

	go func() {
		for range cs.incrCh {
			cs.counter++
			cs.counterCh <- cs.counter
		}
	}()

	return &cs
}

func (c *ChannelCounterService) Close() {
	close(c.counterCh)
	close(c.incrCh)
}

func (c *ChannelCounterService) GetNext() uint64 {
	c.incrCh <- struct{}{}
	return <-c.counterCh
}
