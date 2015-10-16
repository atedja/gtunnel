package tunnel

import (
	"sync"
	"sync/atomic"
)

// Tunnel is a clean wrapper around native Go channel to allow cleanly
// closing the channel without throwing a panic.
// When Tunnel is closed, it waits until all goroutines waiting to
// send data into the tunnel are unblocked.
//
type Tunnel struct {
	wg          *sync.WaitGroup
	mutex       *sync.Mutex
	closed      *atomic.Value
	fullyClosed *atomic.Value
	closingDone chan bool

	// The underlying native Go channel.
	// This is made public only so you can read data from the channel using
	// the standard Select Statement.
	// However, you should always use Send() to send data for better control
	// and throughput.
	Channel chan interface{}
}

// Creates a new Tunnel with no buffer.
//
func NewDefaultTunnel() *Tunnel {
	return &Tunnel{
		wg:          &sync.WaitGroup{},
		mutex:       &sync.Mutex{},
		closed:      &atomic.Value{},
		fullyClosed: &atomic.Value{},
		closingDone: make(chan bool, 1),
		Channel:     make(chan interface{}),
	}
}

// Creates a new Tunnel with buffer.
//
func NewTunnel(buffer int64) *Tunnel {
	return &Tunnel{
		wg:          &sync.WaitGroup{},
		mutex:       &sync.Mutex{},
		closed:      &atomic.Value{},
		fullyClosed: &atomic.Value{},
		closingDone: make(chan bool, 1),
		Channel:     make(chan interface{}, buffer),
	}
}

// Send data into this Tunnel.
// Will yield error if channel is closed.
// You should always use this to send data to channel.
//
func (self *Tunnel) Send(v interface{}) error {
	if self.IsClosed() == true {
		return ErrClosedTunnel
	}
	self.wg.Add(1)
	self.Channel <- v
	self.wg.Done()
	return nil
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
// This method is blocking. After this method completes,
// IsFullyClosed() will always return true, and length of channel equals to 0.
//
func (self *Tunnel) Wait() {
	<-self.closingDone
	return
}

// Close this Tunnel.
//
func (self *Tunnel) Close() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if !self.IsClosed() {
		self.closed.Store(true)
		go func() {
			self.wg.Wait()
			close(self.Channel)
			self.fullyClosed.Store(true)
			self.closingDone <- true
		}()
	}
}
