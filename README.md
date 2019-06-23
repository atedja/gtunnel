# tunnel

[![GoDoc](https://godoc.org/github.com/atedja/gtunnel?status.svg)](https://godoc.org/github.com/atedja/gtunnel)[![Build Status](https://travis-ci.org/atedja/gtunnel.svg?branch=master)](https://travis-ci.org/atedja/gtunnel)

Tunnel is a clean and efficient wrapper around native Go channels that allows you to:
* Close a channel multiple times without throwing a panic.
* Detect if channel has been closed.

### Quick Example

```go
package main

import (
	"fmt"
	"github.com/atedja/gtunnel"
)

func main() {
	tn := tunnel.NewUnbuffered()
	for i := 0; i < 10; i++ {
		go func() {
			err := tn.Send("something")
			if err != nil {
				fmt.Println("Channel is closed!")
			}
		}()
	}

	go func() {
		for data := range tn.Out() {
			fmt.Println(data.(string))
		}
	}()

	tn.Close()
	tn.Wait()

	return
}
```

Read the `tunnel_test.go` for extreme cases.


### Basic Usage

#### Creating a tunnel

    var tn *tunnel.Tunnel
    tn = tunnel.NewBuffered(10)  // for buffered channel
    tn = tunnel.NewUnbuffered()  // for unbuffered

#### Sending data to tunnel

    tn.Send("data")

#### Reading from tunnel

Tunnels use `interface{}` internally, so you need to cast data back to its original format, e.g. `data.(T)`

    for data := range tn.Out() {
        mydata = data.(string)
        ....
    }

#### Close tunnel

    tn.Close()

#### Check if tunnel is closed

    tn.IsClosed()
