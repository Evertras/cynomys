package sender

import (
	"fmt"
	"net"
)

type TCPSender struct {
	broadcastAddr net.TCPAddr
	conn          *net.TCPConn
}

func NewTCPSender(addr net.TCPAddr) *TCPSender {
	return &TCPSender{
		broadcastAddr: addr,
	}
}

func (s *TCPSender) Disconnect() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *TCPSender) Send(data []byte) error {
	if s.conn == nil {
		c, err := net.DialTCP("tcp4", nil, &s.broadcastAddr)

		if err != nil {
			return fmt.Errorf("net.DialTCP: %w", err)
		}

		s.conn = c
	}

	_, err := s.conn.Write(data)

	if err != nil {
		return fmt.Errorf("s.conn.Write: %w", err)
	}

	return nil
}
