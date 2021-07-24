package dependency

import (
	"sync"
)

// Wrapper around WaitGroups to facilitate waiting on dependencies to be ready
type dependency struct {
	// WaitGroup used internally to store and manage the ready state
	waiter sync.WaitGroup
}

// Create new instance of dependency, initially set to not-ready state
func New() dependency {
	var dep dependency
	// Initially increment WaitGroup so that all calls to Await() will lock
	dep.waiter.Add(1)
	return dep
}

// Wait for dependency to be marked as ready
// Should run immediately if dependency is already ready
func (dep *dependency) Await() {
	dep.waiter.Wait()
}

// Mark dependency as ready, run all waiting functions
func (dep *dependency) Ready() {
	// Decrement WaitGroup to 0, permanently allowing waiting functions
	dep.waiter.Done()
}
