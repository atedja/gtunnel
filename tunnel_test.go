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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	tn := NewUnbuffered()
	for i := 0; i < 10; i++ {
		go func() {
			tn.Send("something")
		}()
	}

	go func() {
		for data := range tn.Out() {
			assert.Equal(t, "something", data.(string))
		}
	}()

	tn.Close()
	tn.Wait()

	return
}

func TestTunnelBasicCase(t *testing.T) {
	fmt.Println("Testing basic case")
	th := NewBuffered(4)
	th.Send("yo1")
	th.Send("yo2")
	th.Send("yo3")
	th.Send("yo4")

	var data interface{}
	var ok bool
	data, ok = <-th.Out()
	assert.Equal(t, true, ok)
	assert.Equal(t, "yo1", data.(string))

	data, ok = <-th.Out()
	assert.Equal(t, true, ok)
	assert.Equal(t, "yo2", data.(string))

	data, ok = <-th.Out()
	assert.Equal(t, true, ok)
	assert.Equal(t, "yo3", data.(string))

	data, ok = <-th.Out()
	assert.Equal(t, true, ok)
	assert.Equal(t, "yo4", data.(string))

	// Hammering Close() shouldn't affect anything
	th.Close()
	th.Close()
	th.Close()

	th.Wait()
	th.Wait()
	th.Wait()

	assert.Equal(t, true, th.IsClosed())
	assert.Equal(t, 0, th.Len())
}

func TestTunnelBufferOverflow(t *testing.T) {
	fmt.Println("Testing more data into channel")
	th := NewBuffered(100)

	// Attempt to write more than the tunnel can handle.
	for i := 0; i < 2000; i++ {
		go th.Send(1)
	}

	// Let the loop above run for some time
	time.Sleep(100 * time.Millisecond)

	// Close this tunnel, any further attempt to push should fail.
	th.Close()
	err := th.Send(1)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrClosedTunnel)

	// Flush them out
	go func() {
		for c := range th.Out() {
			assert.Equal(t, 1, c.(int))
		}
	}()

	// Wait() will automatically unblock once flushing completes.
	th.Wait()
	assert.Equal(t, true, th.IsClosed())
	assert.Equal(t, 0, th.Len())

	// Further attempt at reading will fail
	_, ok := <-th.Out()
	assert.Equal(t, false, ok)
}

func TestTunnelGoCrazyUnbuffered(t *testing.T) {
	fmt.Println("Testing lots of goroutines with unbuffered tunnel")
	th := NewUnbuffered()

	// reader
	readerDone := make(chan bool)
	go func() {
		for c := range th.Out() {
			assert.Equal(t, "data", c.(string))
		}
		readerDone <- true
	}()

	// writer
	writerDone := make(chan bool)
	go func() {
		for !th.IsClosed() {
			go th.Send("data")
		}
		writerDone <- true
	}()

	// Wait to make reader and writer go crazy
	time.Sleep(1000 * time.Millisecond)

	th.Close()
	th.Wait()
	assert.Equal(t, true, th.IsClosed())
	assert.Equal(t, 0, th.Len())

	<-readerDone
	<-writerDone
}

func TestTunnelGoCrazyBuffered(t *testing.T) {
	fmt.Println("Testing lots of goroutines with buffered tunnel")
	th := NewBuffered(100)

	// reader
	readerDone := make(chan bool)
	go func() {
		for c := range th.Out() {
			assert.Equal(t, "data", c.(string))
		}
		readerDone <- true
	}()

	// writer
	writerDone := make(chan bool)
	go func() {
		for !th.IsClosed() {
			go th.Send("data")
		}
		writerDone <- true
	}()

	// Wait to make reader and writer go crazy
	time.Sleep(1000 * time.Millisecond)

	th.Close()
	th.Wait()
	assert.Equal(t, true, th.IsClosed())
	assert.Equal(t, 0, th.Len())

	<-readerDone
	<-writerDone
}
