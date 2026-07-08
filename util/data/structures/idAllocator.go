package structures

import (
	"math"
	"sync"
)

const (
	DefaultIdAllocatorCapacity = 1000 // 8KB of memory (1000 * 8 bytes) for buffer
)

// IdAllocator manages a pool of IDs using a thread-safe ring buffer.
type IdAllocator struct {
	mu        *sync.Mutex
	capacity  int
	buffer    []int
	head      int
	tail      int
	allocated int
}

// Allocate retrieves the next available ID. Returns false if the buffer is exhausted.
func (a *IdAllocator) Allocate() (int, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.allocated == a.capacity {
		return -1, false
	}

	id := a.tail
	a.buffer[a.tail] = id

	a.tail = (a.tail + 1) % a.capacity
	a.allocated++

	return id, true
}

// Release returns an ID to the ring buffer so it can be reused later.
func (a *IdAllocator) Release(id int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.allocated == 0 {
		// no ids are allocated, so nothing to release
		return
	}

	a.buffer[a.head] = id

	a.head = (a.head + 1) % a.capacity
	a.allocated--
}

// Used the number of allocated ids
func (a *IdAllocator) Used() int {
	return a.allocated
}

// Free the number of available ids
func (a *IdAllocator) Free() int {
	return a.capacity - a.allocated
}

// NewIdAllocator creates an allocator with a fixed maximum size.
// the buffer will use 8*capacity bytes of memory (e.g. capacity=1000 uses 8KB of memory)
func NewIdAllocator(capacity int) *IdAllocator {
	if capacity <= 0 {
		capacity = DefaultIdAllocatorCapacity
	}

	if capacity > math.MaxInt {
		capacity = math.MaxInt
	}

	return &IdAllocator{
		mu:        &sync.Mutex{},
		capacity:  capacity,
		buffer:    make([]int, capacity),
		head:      0,
		tail:      0,
		allocated: 0,
	}
}
