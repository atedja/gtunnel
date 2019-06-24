package tunnel

import (
	"sync"
)

// Simple Semaphore implementation with buffer and closer.
// Acquire() is guaranteed to succeed as long as buffer is not empty.
// If buffer is empty, Acquire() blocks until either a resource has been released back or Semaphore is closed, which will return an error.
// Closing a Semaphore does not guarantee future Acquire() calls to fail immediately, as it depends on the size of the underlying buffer and the thread scheduling.
// However, once Close() returns, all future calls to Acquire() is guaranteed to fail.
type Semaphore struct {
	sync.Mutex
	cond       *sync.Cond
	buffer     chan byte
	bufferSize int
	closed     bool
	returned   int
	notUsed    int
}

func NewSemaphore(bufferSize int) *Semaphore {
	sem := &Semaphore{
		buffer:     make(chan byte, bufferSize),
		cond:       &sync.Cond{L: &sync.Mutex{}},
		bufferSize: bufferSize,
		closed:     false,
		returned:   0,
		notUsed:    0,
	}

	for i := 0; i < bufferSize; i++ {
		sem.buffer <- 1
	}

	return sem
}

// Non-blocking unless buffer is drained.
// Returns error if buffer is drained and Semaphore has been closed.
// If this function returns nil, it is guaranteed that thread has acquired a resource, such that a call to Wait() will block until the resource is released.
func (sem *Semaphore) Acquire() error {
	if _, ok := <-sem.buffer; !ok {
		return ErrClosedSemaphore
	}
	return nil
}

// Return resource back to Semaphore.
// Do not call Release without calling Acquire first.
func (sem *Semaphore) Release() {
	sem.Lock()
	defer sem.Unlock()
	if !sem.closed {
		sem.buffer <- 1
	} else {
		sem.returned++
		sem.cond.Signal()
	}
}

// Count the number of reclaimed resources
func (sem *Semaphore) count() int {
	sem.Lock()
	defer sem.Unlock()
	return sem.notUsed + sem.returned
}

// Wait blocks until all resources are released back.
func (sem *Semaphore) Wait() {
	sem.cond.L.Lock()
	for sem.count() < sem.bufferSize {
		sem.cond.Wait()
	}
	sem.cond.L.Unlock()
}

// Close this Semaphore, and flush out all unused resources.
// After this function returns, all future calls to Acquire is guaranteed to fail.
func (sem *Semaphore) Close() {
	sem.Lock()
	defer sem.Unlock()
	sem.closed = true
	close(sem.buffer)
	for range sem.buffer {
		sem.notUsed++
	}
}
