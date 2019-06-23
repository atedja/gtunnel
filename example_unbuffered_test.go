package tunnel_test

import (
	"fmt"
	"github.com/atedja/gtunnel"
)

func ExampleNewUnbuffered() {
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
}
