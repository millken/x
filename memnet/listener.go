package memnet

import (
	"errors"
	"net"
	"sync"
)

// ErrListenerClosed indicates that the Listener is already closed.
var ErrListenerClosed = errors.New("Listener is already closed: use of closed network connection")

// Listener provides in-memory dialer<->net.Listener implementation.
//
// It may be used either for fast in-process client<->server communications
// without network stack overhead or for client<->server tests.
type Listener struct {
	lock   sync.Mutex
	closed bool
	conns  chan acceptConn
}

type acceptConn struct {
	conn     net.Conn
	accepted chan struct{}
}

// NewListener returns new in-memory dialer<->net.Listener.
func NewListener() *Listener {
	return &Listener{
		conns: make(chan acceptConn, 1024),
	}
}

// SetLocalAddr sets the (simulated) local address for the listener.
func (ln *Listener) SetLocalAddr(localAddr net.Addr) {
}

// Accept implements net.Listener's Accept.
//
// It is safe calling Accept from concurrently running goroutines.
//
// Accept returns new connection per each Dial call.
func (ln *Listener) Accept() (net.Conn, error) {
	c, ok := <-ln.conns
	if !ok {
		return nil, ErrListenerClosed
	}
	close(c.accepted)
	return c.conn, nil
}

// Close implements net.Listener's Close.
func (ln *Listener) Close() error {
	var err error

	ln.lock.Lock()
	if !ln.closed {
		close(ln.conns)
		ln.closed = true
	} else {
		err = ErrListenerClosed
	}
	ln.lock.Unlock()
	return err
}

type inmemoryAddr int

func (inmemoryAddr) Network() string {
	return "inmemory"
}

func (inmemoryAddr) String() string {
	return "Listener"
}

// Addr implements net.Listener's Addr.
func (ln *Listener) Addr() net.Addr {
	return inmemoryAddr(0)
}

// Dial creates new client<->server connection.
// Just like a real Dial it only returns once the server
// has accepted the connection.
//
// It is safe calling Dial from concurrently running goroutines.
func (ln *Listener) Dial() (net.Conn, error) {
	return ln.DialWithLocalAddr(nil)
}

// DialWithLocalAddr creates new client<->server connection.
// Just like a real Dial it only returns once the server
// has accepted the connection. The local address of the
// client connection can be set with local.
//
// It is safe calling Dial from concurrently running goroutines.
func (ln *Listener) DialWithLocalAddr(local net.Addr) (net.Conn, error) {
	pc := NewPipeConns()

	pc.SetAddresses(local, ln.Addr(), ln.Addr(), local)

	cConn := pc.Conn1()
	sConn := pc.Conn2()
	ln.lock.Lock()
	accepted := make(chan struct{})
	if !ln.closed {
		ln.conns <- acceptConn{sConn, accepted}
		// Wait until the connection has been accepted.
		<-accepted
	} else {
		_ = sConn.Close()
		_ = cConn.Close()
		cConn = nil
	}
	ln.lock.Unlock()

	if cConn == nil {
		return nil, ErrListenerClosed
	}
	return cConn, nil
}
