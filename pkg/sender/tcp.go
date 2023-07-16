package sender

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type TCPSender struct {
	mu sync.RWMutex

	broadcastAddr net.TCPAddr
	conn          *net.TCPConn
	sendInterval  time.Duration
}

func NewTCPSender(addr net.TCPAddr, sendInterval time.Duration) *TCPSender {
	return &TCPSender{
		broadcastAddr: addr,
		sendInterval:  sendInterval,
	}
}

func (s *TCPSender) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *TCPSender) Send(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		c, err := net.DialTCP("tcp4", nil, &s.broadcastAddr)

		if err != nil {
			return fmt.Errorf("net.DialTCP: %w", err)
		}

		s.conn = c
	}

	_, err := s.conn.Write(data)

	if err != nil {
		_ = s.conn.Close()
		s.conn = nil
		return fmt.Errorf("s.conn.Write: %w", err)
	}

	return nil
}

func (s *TCPSender) Run() error {
	s.mu.RLock()
	sendTCPTo := s.broadcastAddr.String()
	sendInterval := s.sendInterval
	s.mu.RUnlock()

	for {
		err := s.Send([]byte("hi"))

		if err != nil {
			log.Printf("Failed to send to %q: %v", sendTCPTo, err)
		}

		time.Sleep(sendInterval)
	}
}
