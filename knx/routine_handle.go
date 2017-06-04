package knx

import (
	"sync"
)

// routineHandle is a handle for a running goroutine.
type routineHandle struct {
	mu     sync.Mutex
	done   chan struct{}
	term   chan struct{}
	closed bool
}

// runRoutine launches a goroutine and creates a handle for it.
func runRoutine(routine func(done <-chan struct{})) *routineHandle {
	handle := &routineHandle{
		done: make(chan struct{}),
		term: make(chan struct{}),
	}

	go func(handle *routineHandle) {
		// Make sure to inform the parent goroutine that we exited.
		defer func(term chan<- struct{}) {
			close(term)
		}(handle.term)

		routine(handle.done)
	}(handle)

	return handle
}

// wait waits until the goroutine has exited.
func (handle *routineHandle) wait() {
	<-handle.term
}

// stop instructs the goroutine to terminate.
func (handle *routineHandle) stop() {
	handle.mu.Lock()
	defer handle.mu.Unlock()

	if !handle.closed {
		close(handle.done)
		handle.closed = true
	}
}
