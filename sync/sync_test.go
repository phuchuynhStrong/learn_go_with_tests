package sync

import (
	"sync"
	"testing"
)

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}

func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		wantedCount := 1000
		counter := NewCounter()

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				counter.Inc()
				wg.Done()
			}()
		}

		wg.Wait()

		assertCounter(t, counter, wantedCount)
	})
}

func NewCounter() *Counter {
	return &Counter{}
}

func assertCounter(t testing.TB, c *Counter, want int) {
	if c.Value() != want {
		t.Errorf("got %v, want %v", c.Value(), want)
	}
}
