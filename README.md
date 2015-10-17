# tunnel

Tunnel is a clean wrapper around native Go channel to allow cleanly closing the
channel without throwing a panic. When Tunnel is closed, it waits until all
goroutines waiting to send data into the tunnel are unblocked.

### [API Documentation](https://godoc.org/github.com/atedja/go-tunnel)

### Why?

Currently in Go, there is no way to know if a channel is closed before sending
data to it. But, if you do send data to a closed channel, it will panic, and
the only way to overcome that is to `recover()` which is Yuck™.

Once closed, Tunnel will reject all new incoming data without panic, and 
your reader can flush out the remaining data in the channel buffer and other 
blocked goroutines are unblocked naturally. You can call `Wait()` to block 
current execution thread until there is no more data in the channel.

See Example below, or read the `tunnel_test.go` for extreme cases.

### Example

```go
package main

import (
	"fmt"
	"github.com/atedja/go-tunnel"
)

func main() {
	tn := tunnel.NewBuffered(1)
	tn.Send("something")

	data := <-tn.Out()
	fmt.Println(data.(string))

	tn.Close()
	err := tn.Send("more")
	if err != nil {
		fmt.Println("channel is closed!")
	}

	tn.Wait()

	return
}
```

