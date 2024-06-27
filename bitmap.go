package main

import (
	"sync/atomic"
)

const (
	ipV4AddressesCount uint32 = (1 << 32) - 1
)

type Bitmap struct {
	bits []uint32
}

// NewBitmap creates a new Bitmap with the given size.
// The size is the number of bits the bitmap should be able to handle.
func NewBitmap(size uint32) *Bitmap {
	// Calculate the number of uint32 needed to store 'size' bits.
	// Each uint32 can store 32 bits.
	numUints := size / 32

	return &Bitmap{
		bits: make([]uint32, numUints),
	}
}

// Set marks the bit at the given position to 1 using atomic operations.
// Returns true if the value was actually updated, and false otherwise.
func (b *Bitmap) Set(bit uint32) bool {
	// Calculate the index in the bits slice.
	index := bit / 32

	// Calculate the position within the uint32 at that index.
	position := bit % 32

	// Load the current value at the index.
	previousValue := atomic.LoadUint32(&b.bits[index])

	// Create a mask with a 1 at the given position.
	mask := uint32(1) << position

	// Compute the new value by ORing the old value with the mask.
	newValue := previousValue | mask

	if previousValue == newValue {
		return false
	}

	// Atomically compare and swap the old value with the new value.
	swapResult := atomic.CompareAndSwapUint32(&b.bits[index], previousValue, newValue)

	// If swap failed, another goroutine already modified the value.
	// So no need to retry the swap.
	return swapResult
}

// IsSet checks if the bit at the given position is 1 using atomic operations.
func (b *Bitmap) IsSet(bit uint32) bool {
	// Calculate the index in the bits slice.
	index := bit / 32

	// Calculate the position within the uint32 at that index.
	position := bit % 32

	// Load the current value at the index.
	value := atomic.LoadUint32(&b.bits[index])

	// Create a mask with a 1 at the given position.
	mask := uint32(1) << position

	// Check the bit using AND operation and compare with the mask.
	return (value & mask) != 0
}
