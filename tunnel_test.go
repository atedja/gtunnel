package tunnel

import (
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
	th := NewBuffered(100)

	// Attempt to write more than the tunnel can handle.
	for i := 0; i < 2000; i++ {
		go th.Send(1)
	}

	// Let the loop above run for some time
	time.Sleep(100 * time.Millisecond)

	// Close this tunnel, any further attempt to push should fail.
	th.Close()

	// Flush them out
	readerDone := make(chan bool)
	go func() {
		for c := range th.Out() {
			assert.Equal(t, 1, c.(int))
		}
		close(readerDone)
	}()

	th.Wait()
	assert.Equal(t, true, th.IsClosed())

	<-readerDone

	// Further attempt at reading will fail
	_, ok := <-th.Out()
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, th.Len())
}

func TestTunnelGoCrazyUnbuffered(t *testing.T) {
	th := NewUnbuffered()

	// reader
	readerDone := make(chan bool)
	go func() {
		for c := range th.Out() {
			assert.Equal(t, "data", c.(string))
		}
		close(readerDone)
	}()

	// writer
	writerDone := make(chan bool)
	go func() {
		var err error
		for err == nil {
			err = th.Send("data")
		}
		close(writerDone)
	}()

	// Wait to make reader and writer go crazy
	time.Sleep(1000 * time.Millisecond)

	th.Close()
	th.Wait()
	assert.Equal(t, true, th.IsClosed())

	<-readerDone
	<-writerDone

	assert.Equal(t, 0, th.Len())
}

func TestTunnelGoCrazyBuffered(t *testing.T) {
	th := NewBuffered(100)

	// reader
	readerDone := make(chan bool)
	go func() {
		for c := range th.Out() {
			assert.Equal(t, "data", c.(string))
		}
		close(readerDone)
	}()

	// writer
	writerDone := make(chan bool)
	go func() {
		var err error
		for err == nil {
			err = th.Send("data")
		}
		close(writerDone)
	}()

	// Wait to make reader and writer go crazy
	time.Sleep(1000 * time.Millisecond)

	th.Close()
	th.Wait()
	assert.Equal(t, true, th.IsClosed())

	<-readerDone
	<-writerDone

	assert.Equal(t, 0, th.Len())
}
