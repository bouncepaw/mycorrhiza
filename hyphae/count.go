package hyphae

import (
	"sync"
)

// Its value is number of all existing hyphae. Hypha mutators are expected to manipulate the value. It is concurrent-safe.
var count = struct {
	value int
	sync.Mutex
}{}

// Set the value of hyphae count to zero.
func ResetCount() {
  count.Lock()
  count.value = 0
  count.Unlock()
}

// Increment the value of hyphae count.
func IncrementCount() {
	count.Lock()
	count.value++
	count.Unlock()
}

// Decrement the value of hyphae count.
func DecrementCount() {
	count.Lock()
	count.value--
	count.Unlock()
}

// Count how many hyphae there are.
func Count() int {
	// it is concurrent-safe to not lock here, right?
	return count.value
}
