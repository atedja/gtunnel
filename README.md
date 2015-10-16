# tunnel

Tunnel is a clean wrapper around native Go channel to allow cleanly closing the
channel without throwing a panic. When Tunnel is closed, it waits until all
goroutines waiting to send data into the tunnel are unblocked.

## [API Documentation](https://godoc.org/github.com/atedja/go-tunnel)

## Example

```go
package main

import (
	"fmt"
	"github.com/atedja/go-tunnel"
)

func main() {
	tn := tunnel.NewTunnel(1)
	tn.Send("something")

	data := <-tn.Channel
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

