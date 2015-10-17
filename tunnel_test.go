package tunnel

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTunnel(t *testing.T) {
	th := NewTunnel(4)
	th.Send("yo1")
	th.Send("yo2")
	th.Send("yo3")
	th.Send("yo4")
	assert.Equal(t, 4, th.Len())

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

	th.Close()
	assert.Equal(t, true, th.IsClosed())
	th.Wait()
	assert.Equal(t, true, th.IsFullyClosed())
	assert.Equal(t, 0, th.Len())
}

func TestTunnelCorrectBuffer(t *testing.T) {
	th := NewTunnel(100)

	// writes 200 times to tunnel that only fits 100
	writeDone := make(chan bool)
	go func() {
		for i := 0; i < 200; i++ {
			go func() {
				th.Send(1)
			}()
		}
		writeDone <- true
	}()
	<-writeDone

	// Close this tunnel, any further attempt to push should fail
	th.Close()
	assert.Equal(t, true, th.IsClosed())
	err := th.Send(1)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrClosedTunnel)

	// Flush them out, there should be 200 of them
	go func() {
		counter := 0
		for c := range th.Out() {
			assert.Equal(t, 1, c.(int))
			counter++
		}
		assert.Equal(t, 200, counter)
	}()

	// Wait() will automatically unblock once flushing completes.
	th.Wait()
	assert.Equal(t, true, th.IsFullyClosed())
	assert.Equal(t, 0, th.Len())

	// Further attempt at reading will fail
	_, ok := <-th.Out()
	assert.Equal(t, false, ok)
}

func TestTunnelGoCrazy(t *testing.T) {
	th := NewDefaultTunnel()

	// reader
	readerDone := make(chan bool)
	go func() {
	MainLoop:
		for {
			select {
			case _, ok := <-th.Out():
				if !ok {
					break MainLoop
				}
			}
		}
		readerDone <- true
	}()

	// writer
	writerDone := make(chan bool)
	go func() {
		for {
			if th.IsClosed() {
				break
			}

			go func() {
				th.Send(1)
			}()
		}
		writerDone <- true
	}()

	// Wait to make reader and writer go crazy
	time.Sleep(1000 * time.Millisecond)

	th.Close()
	assert.Equal(t, true, th.IsClosed())
	th.Wait()
	assert.Equal(t, true, th.IsFullyClosed())
	assert.Equal(t, 0, th.Len())

	<-readerDone
	<-writerDone
}
