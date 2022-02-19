package hyphae

import (
	"sync"
)

// Its value is number of all existing hyphae. NonEmptyHypha mutators are expected to manipulate the value. It is concurrent-safe.
var count = struct {
	value int
	sync.Mutex
}{}

// ResetCount sets the value of hyphae count to zero. Use when reloading hyphae.
func ResetCount() {
	count.Lock()
	count.value = 0
	count.Unlock()
}

// Count how many hyphae there are. This is a O(1), the number of hyphae is stored in memory.
func Count() int {
	count.Lock()
	defer count.Unlock()
	return count.value
}

// incrementCount increments the value of the hyphae counter. Use when creating new hyphae or loading hyphae from disk.
func incrementCount() {
	count.Lock()
	count.value++
	count.Unlock()
}

// decrementCount decrements the value of the hyphae counter. Use when deleting existing hyphae.
func decrementCount() {
	count.Lock()
	count.value--
	count.Unlock()
}
