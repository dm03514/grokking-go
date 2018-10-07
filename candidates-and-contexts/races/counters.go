package races

import "sync"

type Counter struct {
	count int
}

func (c *Counter) Inc() {
	c.count++
}

func (c *Counter) Value() int {
	return c.count
}

func (c *Counter) Set(v int) {
	c.count = v
}

type SynchronizedCounter struct {
	mu    *sync.Mutex
	count int
}

func (c *SynchronizedCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count++
}

func (c *SynchronizedCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}

func (c *SynchronizedCounter) Set(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count = v
}

type MisSynchronizedCounter struct {
	mu    sync.Mutex
	count int
}

func (c MisSynchronizedCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count++
}

func (c *MisSynchronizedCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}

func (c *MisSynchronizedCounter) Set(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count = v
}
