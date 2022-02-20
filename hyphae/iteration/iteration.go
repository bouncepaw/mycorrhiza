// Package iteration provides a handy API for making multiple checks on all hyphae in one go.
package iteration

import (
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"sync"
)

// Iteration represents an iteration over all existing hyphae in the storage. Iteration is done on all existing hyphae. The order of hyphae is not specified. For all hyphae, checks are made.
type Iteration struct {
	sync.Mutex
	checks []func(h hyphae.Hypha) CheckResult
}

// NewIteration constructs an iteration without checks.
func NewIteration() *Iteration {
	return &Iteration{
		checks: make([]func(h hyphae.Hypha) CheckResult, 0),
	}
}

// AddCheck adds the check to the iteration. It is concurrent-safe. Checks are meant to have side-effects.
func (i7n *Iteration) AddCheck(check func(h hyphae.Hypha) CheckResult) {
	i7n.Lock()
	i7n.checks = append(i7n.checks, check)
	i7n.Unlock()
}

func (i7n *Iteration) removeCheck(i int) {
	i7n.checks[i] = i7n.checks[len(i7n.checks)-1]
	i7n.checks = i7n.checks[:len(i7n.checks)-1]
}

// Ignite does the iteration by walking over all hyphae yielded by the iterator used and calling all checks on the hypha. Ignited iterations are not concurrent-safe.
//
// After ignition, you should not use the same Iteration again.
func (i7n *Iteration) Ignite() {
	for h := range hyphae.YieldExistingHyphae() {
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
