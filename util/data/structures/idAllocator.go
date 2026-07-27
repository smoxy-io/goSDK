package structures

import (
	"math"
	"sync"
)

const (
	DefaultIdAllocatorCapacity = 1000 // 1KB of memory (1000 * 1 byte) for buffer
)

// IdAllocator manages a pool of IDs using a thread-safe ring buffer.
type IdAllocator struct {
	mu        *sync.Mutex
	capacity  int
	buffer    []bool
	allocated int
	next      int
}

// Allocate retrieves the next available ID. Returns false if the buffer is exhausted.
func (a *IdAllocator) Allocate() (int, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.allocated == a.capacity || a.next == -1 {
		return -1, false
	}

	defer a.updateNext()

	next := a.next

	a.buffer[next] = true
	a.allocated++

	return next, true
}

// Release returns an ID to the ring buffer so it can be reused later.
func (a *IdAllocator) Release(id int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.allocated == 0 {
		// no ids are allocated, so nothing to release
		return
	}

	// make sure that the id is valid
	if id < 0 || id >= a.capacity {
		return
	}

	if !a.buffer[id] {
		// id is not allocated, so nothing to release
		return
	}

	defer a.updateNext()

	a.buffer[id] = false
	a.allocated--

	if a.next == -1 {
		// all ids were allocated prior to this release. set next to this id
		a.next = id
	}
}

// Used the number of allocated ids
func (a *IdAllocator) Used() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.allocated
}

// Free the number of available ids
func (a *IdAllocator) Free() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.capacity - a.allocated
}

// MUST be called while holding the lock
func (a *IdAllocator) updateNext() {
	if a.allocated == a.capacity {
		a.next = -1
		return
	}

	next := a.next

	if !a.buffer[next] {
		// current next is still available
		// no need to change next pointer
		if next != a.next {
			a.next = next
		}

		return
	}

	// find the next available id
	next++

	for {
		if next == a.capacity {
			// wrap around to the beginning
			next = 0
		}

		if next == a.next {
			// back where we started. no available ids
			a.allocated = a.capacity // updating allocated as it should have prevented us from getting to here
			a.next = -1
			return
		}

		if !a.buffer[next] {
			// this is the next available id
			a.next = next
			return
		}

		next++
	}
}

// NewIdAllocator creates an allocator with a fixed maximum size.
// the buffer will use capacity bytes of memory (e.g. capacity=1000 uses 1KB of memory)
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
		buffer:    make([]bool, capacity),
		next:      0,
		allocated: 0,
	}
}
