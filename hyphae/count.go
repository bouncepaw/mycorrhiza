package hyphae

import (
	"sync"
)

// Its value is number of all existing hyphae. Hypha mutators are expected to manipulate the value. It is concurrent-safe.
var count = struct {
	value int
	sync.Mutex
}{}

// Set the value of hyphae count to zero. Use when reloading hyphae.
func ResetCount() {
	count.Lock()
	count.value = 0
	count.Unlock()
}

// Increment the value of the hyphae counter. Use when creating new hyphae or loading hyphae from disk.
func IncrementCount() {
	count.Lock()
	count.value++
	count.Unlock()
}

// Decrement the value of the hyphae counter. Use when deleting existing hyphae.
func DecrementCount() {
	count.Lock()
	count.value--
	count.Unlock()
}

// Count how many hyphae there are.
func Count() int {
	return count.value
}
