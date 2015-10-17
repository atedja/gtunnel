package tunnel

import (
	//"fmt"
	"sync"
	"sync/atomic"
)

// Tunnel is a clean wrapper around native Go channel to allow cleanly
// closing the channel without throwing a panic.
// When Tunnel is closed, it waits until all goroutines waiting to
// send data into the tunnel are unblocked.
//
type Tunnel struct {
	counter     int64
	cond        *sync.Cond
	mutex       *sync.Mutex
	closed      *atomic.Value
	fullyClosed *atomic.Value
	closingDone chan bool
	channel     chan interface{}
}

// Creates a new Tunnel with no buffer.
//
func NewUnbuffered() *Tunnel {
	tn := &Tunnel{
		counter:     0,
		cond:        &sync.Cond{L: &sync.Mutex{}},
		mutex:       &sync.Mutex{},
		closed:      &atomic.Value{},
		fullyClosed: &atomic.Value{},
		closingDone: make(chan bool, 1),
		channel:     make(chan interface{}),
	}
	return tn
}

// Creates a new Tunnel with buffer.
//
func NewBuffered(buffer int64) *Tunnel {
	tn := &Tunnel{
		counter:     0,
		cond:        &sync.Cond{L: &sync.Mutex{}},
		mutex:       &sync.Mutex{},
		closed:      &atomic.Value{},
		fullyClosed: &atomic.Value{},
		closingDone: make(chan bool, 1),
		channel:     make(chan interface{}, buffer),
	}
	return tn
}

// Send data into this Tunnel.
// Will yield error if channel is closed rather than panics.
// You should always use this to send data to channel.
//
func (self *Tunnel) Send(v interface{}) error {
	if self.IsClosed() {
		return ErrClosedTunnel
	}

	atomic.AddInt64(&self.counter, 1)
	self.cond.L.Lock()
	defer func() {
		atomic.AddInt64(&self.counter, -1)
		self.cond.L.Unlock()
		self.cond.Signal()
	}()
	if self.IsClosed() {
		return ErrClosedTunnel
	}
	self.channel <- v

	return nil
}

func (self *Tunnel) Out() <-chan interface{} {
	return self.channel
}

func (self *Tunnel) Len() int {
	return len(self.channel)
}

// Check if this Tunnel is closed.
//
func (self *Tunnel) IsClosed() bool {
	if self.closed.Load() == true {
		return true
	}
	return false
}

// Check if this Tunnel is fully closed.
// Fully closed is defined by the closure of the underlying channel, length
// of channel is 0, and no further reading/writing is possible.
//
func (self *Tunnel) IsFullyClosed() bool {
	if self.fullyClosed.Load() == true {
		return true
	}
	return false
}

// Wait until this Tunnel is fully closed.
// This method is blocking until all contents are flushed out of the channel.
// After this method completes, IsFullyClosed() will always return true, and
// length of channel equals to 0.
//
func (self *Tunnel) Wait() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if !self.IsFullyClosed() {
		<-self.closingDone
		self.fullyClosed.Store(true)
	}
}

// Close this Tunnel.
//
func (self *Tunnel) Close() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if !self.IsClosed() {
		self.closed.Store(true)
		go func() {
			self.cond.L.Lock()
			for atomic.LoadInt64(&self.counter) != 0 {
				self.cond.Wait()
			}
			self.cond.L.Unlock()
			close(self.channel)
			self.closingDone <- true
			close(self.closingDone)
		}()
	}
}
