package structures

import (
	"testing"
)

func Test_IdAllocator(t *testing.T) {
	pool := NewIdAllocator(-5)

	if pool.Free() != DefaultIdAllocatorCapacity {
		t.Errorf("Expected free capacity to be %d, got %d", DefaultIdAllocatorCapacity, pool.Free())
	}

	pool = NewIdAllocator(5)

	// Allocate 5 IDs
	for i := 0; i < 5; i++ {
		_, ok := pool.Allocate()

		if !ok {
			t.Errorf("Failed to allocate ID %d", i)
		}
	}

	// Try to allocate a 6th ID (Should fail since capacity is 5)
	if _, ok := pool.Allocate(); ok {
		t.Errorf("allocating 6th ID with capacity 5 succeeded. expected failure")
	}

	// Release 2 IDs back into the ring buffer
	pool.Release(0)
	pool.Release(1)

	// Allocate again to verify the recycled IDs are re-issued
	for i := 0; i < 2; i++ {
		if id, ok := pool.Allocate(); !ok {
			t.Errorf("Failed to re-allocate ID %d", i)
		} else {
			if id != 0 && id != 1 {
				t.Errorf("invalid re-allocated ID %d. expected 0 or 1", id)
			}
		}
	}
}
