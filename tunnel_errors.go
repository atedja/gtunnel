package tunnel

import (
	"errors"
)

var ErrClosedTunnel = errors.New("Tunnel is closed")
