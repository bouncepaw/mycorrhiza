package hyphae

import (
	"sync"
)

// Iteration represents an iteration over all hyphae in the storage. You may use it instead of directly iterating using hyphae.YieldExistingHyphae when you want to do n checks at once instead of iterating n times.
type Iteration struct {
	sync.Mutex
	iterator func() chan Hypher
	checks   []func(h Hypher) CheckResult
}

// NewIteration constructs an iteration without checks.
func NewIteration() *Iteration {
	return &Iteration{
		iterator: YieldExistingHyphae,
		checks:   make([]func(h Hypher) CheckResult, 0),
	}
}

// AddCheck adds the check to the iteration. It is concurrent-safe.
func (i7n *Iteration) AddCheck(check func(h Hypher) CheckResult) {
	i7n.Lock()
	i7n.checks = append(i7n.checks, check)
	i7n.Unlock()
}

func (i7n *Iteration) removeCheck(i int) {
	i7n.checks[i] = i7n.checks[len(i7n.checks)-1]
	i7n.checks = i7n.checks[:len(i7n.checks)-1]
}

// Ignite does the iteration by walking over all hyphae yielded by the iterator used and calling all checks on the hypha. Ignited iterations are not concurrent-safe.
func (i7n *Iteration) Ignite() {
	for h := range i7n.iterator() {
		for i, check := range i7n.checks {
			if res := check(h); res == CheckForgetMe {
				i7n.removeCheck(i)
			}
		}
	}
}

// CheckResult is a result of an iteration check.
type CheckResult int

const (
	// CheckContinue is returned when the check wants to be used next time too.
	CheckContinue CheckResult = iota
	// CheckForgetMe is returned when the check wants to be forgotten and not used anymore.
	CheckForgetMe
)
