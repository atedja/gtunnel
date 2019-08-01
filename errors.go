package tunnel

import (
	"errors"
)

// ErrClosedTunnel is returned when tunnel is already closed.
var ErrClosedTunnel = errors.New("gtunnel tunnel is closed")

// ErrClosedSemaphore is returned when semaphore is already closed.
var ErrClosedSemaphore = errors.New("gtunnel semaphore is closed")
