package sender

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type UDPSender struct {
	mu sync.RWMutex

	broadcastAddr net.UDPAddr
	conn          *net.UDPConn
	sendInterval  time.Duration
}

func NewUDPSender(addr net.UDPAddr, sendInterval time.Duration) *UDPSender {
	return &UDPSender{
		broadcastAddr: addr,
		sendInterval:  sendInterval,
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

func (s *UDPSender) Run() error {
	s.mu.RLock()
	sendUDPTo := s.broadcastAddr.String()
	sendInterval := s.sendInterval
	s.mu.RUnlock()

	for {
		err := s.Send([]byte("hi"))
		if err != nil {
			log.Printf("Failed to send to %q: %v", sendUDPTo, err)
		}
		time.Sleep(sendInterval)
	}
}
