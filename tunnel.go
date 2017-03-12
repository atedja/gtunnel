/*
Copyright 2015-2017 Albert Tedja

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	coffee      *sync.Once
	closed      *atomic.Value
	closingDone chan bool
	channel     chan interface{}
	semaphore   *Semaphore
}

// Creates a new Tunnel with no buffer.
//
func NewUnbuffered() *Tunnel {
	tn := &Tunnel{
		coffee:      &sync.Once{},
		closed:      &atomic.Value{},
		closingDone: make(chan bool, 1),
		channel:     make(chan interface{}),
		semaphore:   NewSemaphore(1),
	}
	return tn
}

// Creates a new Tunnel with buffer.
//
func NewBuffered(buffer int) *Tunnel {
	tn := &Tunnel{
		coffee:      &sync.Once{},
		closed:      &atomic.Value{},
		closingDone: make(chan bool, 1),
		channel:     make(chan interface{}, buffer),
		semaphore:   NewSemaphore(buffer),
	}
	return tn
}

// Send data into this Tunnel.
// Will yield error if channel is closed rather than panics.
// You should always use this to send data to channel.
//
func (self *Tunnel) Send(v interface{}) error {
	if self.semaphore.Acquire() != nil {
		return ErrClosedTunnel
	}
	defer self.semaphore.Release()

	self.channel <- v

	return nil
}

func (self *Tunnel) Out() <-chan interface{} {
	return self.channel
}

func (self *Tunnel) Len() int {
	return len(self.channel)
}

func (self *Tunnel) IsClosed() bool {
	if self.closed.Load() == true {
		return true
	}
	return false
}

// Wait until closing process completes.
//
func (self *Tunnel) Wait() {
	<-self.closingDone
}

// Close this Tunnel. Always Be Closing.
//
func (self *Tunnel) Close() {
	// coffee is for closers only.
	go func() {
		self.coffee.Do(func() {
			self.semaphore.Close()
			self.semaphore.Wait()
			close(self.channel)
			self.closed.Store(true)
			self.closingDone <- true
			close(self.closingDone)
		})
	}()
}
