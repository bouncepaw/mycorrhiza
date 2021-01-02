package hypha

import (
	"sync"
)

type count struct {
	value uint
	sync.Mutex
}

// Count is a global variable. Its value is number of all existing hyphae. Hypha mutators are expected to manipulate the value. It is concurrent-safe.
var Count = count{}

// Increment the value of Count.
func (c *count) Increment() {
	c.Lock()
	c.value++
	c.Unlock()
}

// Decrement the value of Count.
func (c *count) Decrement() {
	c.Lock()
	c.value--
	c.Unlock()
}

// Get value of Count.
func (c *count) Value() uint {
	// it is concurrent-safe to not lock here, right?
	return c.value
}
