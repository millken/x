package udp

import (
	"bufio"
	"net"
)

type UDPConfig struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

// UDPWriter is a writer that sends messages to a remote UDP server.
type UDPWriter struct {
	cfg          UDPConfig
	bufferWriter *bufio.Writer
}

// NewUDPWriter creates a new UDPWriter.
func NewUDPWriter(cfg UDPConfig) (*UDPWriter, error) {
	wr := &UDPWriter{
		cfg: cfg,
	}
	udpAddr, err := net.ResolveUDPAddr(cfg.Network, cfg.Address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP(cfg.Network, nil, udpAddr)
	if err != nil {
		return nil, err
	}
	wr.bufferWriter = bufio.NewWriter(conn)
	return wr, nil
}

// Write writes the message to the network connection.
func (w *UDPWriter) Write(p []byte) (n int, err error) {
	return w.bufferWriter.Write(p)
}
