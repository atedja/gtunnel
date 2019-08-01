package tunnel

import (
	"sync"
	"sync/atomic"
)

// Tunnel is a clean wrapper around native Go channel to allow cleanly closing the channel without throwing a panic.
// When Tunnel is closed, it waits until all goroutines waiting to send data into the tunnel are unblocked.
type Tunnel struct {
	once        *sync.Once
	closed      *atomic.Value
	closingDone chan bool
	channel     chan interface{}
	semaphore   *Semaphore
}

// NewUnbuffered creates a new Tunnel with no buffer.
func NewUnbuffered() *Tunnel {
	tn := &Tunnel{
		once:        &sync.Once{},
		closed:      &atomic.Value{},
		closingDone: make(chan bool),
		channel:     make(chan interface{}),
		semaphore:   NewSemaphore(1),
	}
	return tn
}

// NewBuffered creates a new Tunnel with buffer.
func NewBuffered(buffer int) *Tunnel {
	tn := &Tunnel{
		once:        &sync.Once{},
		closed:      &atomic.Value{},
		closingDone: make(chan bool),
		channel:     make(chan interface{}, buffer),
		semaphore:   NewSemaphore(buffer),
	}
	return tn
}

// Sends data into this Tunnel.
// Will yield error if channel is closed rather than panics.
// You should always use this to send data to channel.
func (tun *Tunnel) Send(v interface{}) error {
	if tun.semaphore.Acquire() != nil {
		return ErrClosedTunnel
	}
	defer tun.semaphore.Release()

	tun.channel <- v

	return nil
}

func (tun *Tunnel) Out() <-chan interface{} {
	return tun.channel
}

func (tun *Tunnel) Len() int {
	return len(tun.channel)
}

func (tun *Tunnel) IsClosed() bool {
	if tun.closed.Load() == true {
		return true
	}
	return false
}

// Wait until closing process completes.
func (tun *Tunnel) Wait() {
	<-tun.closingDone
}

// Close this Tunnel.
func (tun *Tunnel) Close() {
	go func() {
		tun.once.Do(func() {
			tun.semaphore.Close()
			tun.semaphore.Wait()
			close(tun.channel)
			tun.closed.Store(true)
			close(tun.closingDone)
		})
	}()
}
