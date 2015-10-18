package tunnel

import (
	"errors"
)

var ErrClosedTunnel = errors.New("Tunnel is closed")
var ErrClosedSemaphore = errors.New("Semaphore is closed")
