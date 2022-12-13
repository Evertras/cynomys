package sender

import (
	"fmt"
	"net"
)

type UDPSender struct {
	broadcastAddr net.UDPAddr
	conn          *net.UDPConn
}

func NewUDPSender(addr net.UDPAddr) *UDPSender {
	return &UDPSender{
		broadcastAddr: addr,
	}
}

func (s *UDPSender) Send(data []byte) error {
	if s.conn == nil {
		c, err := net.DialUDP("udp4", nil, &s.broadcastAddr)

		if err != nil {
			return fmt.Errorf("net.DialUDP: %w", err)
		}

		s.conn = c
	}

	_, err := s.conn.Write(data)

	if err != nil {
		return fmt.Errorf("s.conn.Write: %w", err)
	}

	return nil
}
