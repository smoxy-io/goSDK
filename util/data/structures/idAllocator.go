package structures

import "sync"

// IdAllocator manages a pool of IDs using a thread-safe ring buffer.
type IdAllocator[T int | uint | int64 | uint64 | int32 | uint32] struct {
	mu        *sync.Mutex
	capacity  T
	buffer    []T
	head      T
	tail      T
	allocated T
}

// Allocate retrieves the next available ID. Returns false if the buffer is exhausted.
func (a *IdAllocator[T]) Allocate() (T, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.allocated == a.capacity {
		return 0, false
	}

	id := a.tail
	a.buffer[a.tail] = id

	a.tail = (a.tail + 1) % a.capacity
	a.allocated++

	return id, true
}

// Release returns an ID to the ring buffer so it can be reused later.
func (a *IdAllocator[T]) Release(id T) {
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

func (a *IdAllocator[T]) Used() T {
	return a.allocated
}

func (a *IdAllocator[T]) Free() T {
	return a.capacity - a.allocated
}

// NewIdAllocator creates an allocator with a fixed maximum size.
func NewIdAllocator[T int | uint | int64 | uint64 | int32 | uint32](capacity T) *IdAllocator[T] {
	if capacity <= 0 {
		capacity = 1024
	}

	return &IdAllocator[T]{
		mu:        &sync.Mutex{},
		capacity:  capacity,
		buffer:    make([]T, capacity),
		head:      0,
		tail:      0,
		allocated: 0,
	}
}
