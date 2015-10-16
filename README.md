# tunnel
--
    import "github.com/atedja/go-tunnel"


## Usage

```go
var ErrClosedTunnel = errors.New("Tunnel is closed")
```

#### type Tunnel

```go
type Tunnel struct {

	// The underlying native Go channel.
	// This is made public only so you can read data from the channel using
	// the standard Select Statement.
	// However, you should always use Send() to send data for better control
	// and throughput.
	Channel chan interface{}
}
```

Tunnel is a clean wrapper around native Go channel to allow cleanly closing the
channel without throwing a panic. When Tunnel is closed, it waits until all
goroutines waiting to send data into the tunnel are unblocked.

#### func  NewDefaultTunnel

```go
func NewDefaultTunnel() *Tunnel
```
Creates a new Tunnel with no buffer.

#### func  NewTunnel

```go
func NewTunnel(buffer int64) *Tunnel
```
Creates a new Tunnel with buffer.

#### func (*Tunnel) Close

```go
func (self *Tunnel) Close()
```
Close this Tunnel.

#### func (*Tunnel) IsClosed

```go
func (self *Tunnel) IsClosed() bool
```
Check if this Tunnel is closed.

#### func (*Tunnel) IsFullyClosed

```go
func (self *Tunnel) IsFullyClosed() bool
```
Check if this Tunnel is fully closed. Fully closed is defined by the closure of
the underlying channel, length of channel is 0, and no further reading/writing
is possible.

#### func (*Tunnel) Send

```go
func (self *Tunnel) Send(v interface{}) error
```
Send data into this Tunnel. Will yield error if channel is closed. You should
always use this to send data to channel.

#### func (*Tunnel) Wait

```go
func (self *Tunnel) Wait()
```
Wait until this Tunnel is fully closed. This method is blocking. After this
method completes, IsFullyClosed() will always return true, and length of channel
equals to 0.
