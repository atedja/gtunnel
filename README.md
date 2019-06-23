# tunnel

[![GoDoc](https://godoc.org/github.com/atedja/gtunnel?status.svg)](https://godoc.org/github.com/atedja/gtunnel) [![Build Status](https://travis-ci.org/atedja/gtunnel.svg?branch=master)](https://travis-ci.org/atedja/gtunnel)

Tunnel is a clean and efficient wrapper around native Go channels that allows you to:
* Close a channel multiple times without throwing a panic.
* Detect if channel has been closed.
