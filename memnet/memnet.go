package memnet

import (
	"context"
	"errors"
	"net"
	"sync"
)

var (
	_memnet = sync.Map{}
)

func Listen(network, _ string) (net.Listener, error) {
	var ln net.Listener
	if ln, ok := _memnet.Load(network); ok {
		return ln.(*Listener), nil

	}
	ln = &Listener{
		conns: make(chan acceptConn, 1024),
	}
	_memnet.Store(network, ln)

	return ln, nil
}

func DialContext(ctx context.Context, network, _ string) (net.Conn, error) {
	ln, ok := _memnet.Load(network)
	if !ok {
		return nil, errors.New("memnet: DialContext called for unknown network " + network)
	}
	return ln.(*Listener).Dial()
}
