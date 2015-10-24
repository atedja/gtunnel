# tunnel

Tunnel is a clean and efficient wrapper around native Go channels that allows you to:
* Cleanly close a channel without throwing a panic.
* Detect if it has been closed.

### [API Documentation](https://godoc.org/github.com/atedja/go-tunnel)

### Usage

Currently in Go, there is no way to know if a channel is closed before sending data to it. But, if you do send data to a closed channel, it will panic, and the only way to recover is to `recover()` which is Yuck™.

With tunnels, you can close it by calling `Close()`, and detect if it has been closed by calling `IsClosed()`. It is safe to call `Close()` multiple times from multiple places and multiple goroutines, so you can be at ease and focus on the communication.

After calling `Close()`, tunnel does not immediately close the underlying channel. It blocks future attempts to send data by returning an error, waits for all pending goroutines to finish sending their data, then safely and automatically closes it.  You may optionally call `Wait()` to block the current execution thread until tunnel is fully closed (although data may still be present in the underlying channel if you use buffered channels).

See Example below, or read the `tunnel_test.go` for extreme cases.


### Replacing Channels with Tunnels

Given type `T` and buffer length `L`:

* For unbuffered channels, replace all your `c := make(chan T)` with `c := tunnel.NewUnbuffered()`
* For buffered channels, replace all your `c := make(chan T, L)` with `c := tunnel.NewBuffered(L)`
* Replace all your `c <- data` with `c.Send(data)`
* Replace all your `data, ok := <-c` with `data, ok := <-c.Out()`. Tunnels use `interface{}` internally, so you need to cast data back to its original format, e.g. `data.(T)`
* Close it by calling `c.Close()`
* Check if it's closed by calling `c.IsClosed()`


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
