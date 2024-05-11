package udp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func testListener(t *testing.T) (*net.UDPConn, func()) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:1234")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return conn, func() {
		conn.Close()
	}
}

func TestUDPWriter_Write(t *testing.T) {
	r := require.New(t)
	// Create a new UDPListener
	conn, cleanup := testListener(t)
	defer cleanup()

	// Read a message from the network connection
	go func() {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			t.Fatal(err)
		}
		r.Equal("test1", string(buf[:n]))
	}()
	// Create a new UDPWriter
	w, err := NewUDPWriter(UDPConfig{
		Network: "udp",
		Address: "localhost:1234",
	})
	r.NoError(err)
	// Write a message to the network connection
	_, err = w.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	r.NoError(err)
}
