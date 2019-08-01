# tunnel

[![GoDoc](https://godoc.org/github.com/atedja/gtunnel?status.svg)](https://godoc.org/github.com/atedja/gtunnel) [![Build Status](https://travis-ci.org/atedja/gtunnel.svg?branch=master)](https://travis-ci.org/atedja/gtunnel)

Tunnel is a clean and efficient wrapper around native Go channels that allows you to:
* Gracefully close a channel multiple times without throwing a panic.
* Get an error when attempting to send data to a closed channel.
