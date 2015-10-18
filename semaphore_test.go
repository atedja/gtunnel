package tunnel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	fmt.Println("Testing Semaphore basic functions")

	sm := NewSemaphore(10)
	for i := 0; i < 5; i++ {
		assert.Nil(t, sm.Acquire())
	}
	sm.Close()

	// Only half remains in the buffer
	assert.Equal(t, 5, sm.Count())
	for i := 0; i < 5; i++ {
		sm.Release()
	}

	assert.Equal(t, 10, sm.Count())
	sm.Wait()
}

func TestSemaphoreBufferOne(t *testing.T) {
	fmt.Println("Testing Semaphore of size one")

	sm := NewSemaphore(1)
	assert.Nil(t, sm.Acquire())
	sm.Close()
	err := sm.Acquire()
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrClosedSemaphore)
	sm.Release()
	sm.Wait()
	assert.Equal(t, 1, sm.Count())
}

func TestSemaphoreAcquiredAll(t *testing.T) {
	fmt.Println("Testing Semaphore all resources are acquired")

	sm := NewSemaphore(10)
	for i := 0; i < 10; i++ {
		assert.Nil(t, sm.Acquire())
	}
	sm.Close()
	err := sm.Acquire()
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrClosedSemaphore)

	go func() {
		for i := 0; i < 10; i++ {
			sm.Release()
		}
	}()
	sm.Wait()
}

func TestSemaphoreThreaded(t *testing.T) {
	fmt.Println("Testing Semaphore in threaded situations")

	// This test should take roughly 1 second because we have 10 resources
	// in the semaphore, and 100 goroutines attempting to acquire those
	// resources, and each goroutine takes 100 ms to complete.

	sm := NewSemaphore(10)
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			if sm.Acquire() != nil {
				return
			}
			defer sm.Release()

			// some long operation
			time.Sleep(100 * time.Millisecond)
		}()
	}
	wg.Wait()
	sm.Close()
	sm.Wait()
	assert.Equal(t, 10, sm.Count())
}
